package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/LiangXianSen/go-utils/errors"
)

type Logger struct {
	pipeline logrus.Fields
	*logrus.Logger
}

// Records writes content into logger Fields
func (log *Logger) Records(msg string, fd interface{}) {
	log.pipeline[msg] = fd
}

// NewLogger returns logrus instance
func NewLogger() *Logger {
	return &Logger{
		pipeline: logrus.Fields{},
		Logger:   logrus.New(),
	}
}

type LoggerOptions struct {
	Application  string
	Version      string
	EnableOutput bool
	EnableDebug  bool
}

// Logger is a middleware which provide a logger in ctx.
// Records each handling on os.stdout.
// nolint:funlen
func LoggerM(opt LoggerOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		logS := NewLogger()
		logS.SetFormatter(&logrus.JSONFormatter{})
		logS.SetOutput(os.Stdout)
		logS.SetLevel(logrus.InfoLevel)

		if opt.EnableDebug {
			logS.SetLevel(logrus.DebugLevel)
		}
		c.Set("logger", logS)

		info := logrus.Fields{
			"start":       start,
			"path":        path,
			"method":      method,
			"client_ip":   clientIP,
			"version":     opt.Version,
			"application": opt.Application,
		}

		// Records request parameters
		params := c.Request.URL.Query()
		if len(params) != 0 {
			info["params"] = params
		}

		// Records request body
		requestBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}
		if len(requestBody) != 0 {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
			if c.Request.Header.Get("Content-Type") == "application/json" {
				info["request_body"] = json.RawMessage(requestBody)
			} else {
				info["request_body"] = string(requestBody)
			}
		}

		// Replace gin writer for backup writer stream
		writer := new(multiWriter)
		writer.ctx = c
		writer.ResponseWriter = c.Writer
		c.Writer = writer

		c.Next()

		statusCode := c.Writer.Status()
		requestID := c.GetString("request_id")
		duration := Milliseconds(time.Since(start))
		info["status_code"] = statusCode
		info["request_id"] = requestID
		info["runtime"] = duration

		// Get response from multiWriter
		resp, _ := c.Get("response")
		if buf, ok := resp.(map[string]interface{}); ok {
			info["response"] = buf
		} else {
			info["response"] = resp
		}

		// Writes pipeline from handlers
		pipeline := make(map[string]interface{})
		for k, v := range logS.pipeline {
			pipeline[k] = v
		}
		if len(pipeline) != 0 {
			info["pipeline"] = pipeline
		}

		filterBodyTooLong(info)

		if err, ok := c.Get("error"); ok {
			info["error"] = fmt.Sprintf("%v", err)
			if opt.EnableDebug {
				if e, ok := err.(*errors.Error); ok && e.Stack() != nil {
					info["error"] = fmt.Sprintf("%+v", err.(*errors.Error).Stack())
				}
			}
			logS.WithFields(info).Error("error occurred")
			return
		}

		if opt.EnableOutput {
			logS.WithFields(info).Info("finished")
		}
	}
}

// multiWriter is a backup of gin responseWriter
type multiWriter struct {
	gin.ResponseWriter
	ctx *gin.Context
}

func (w *multiWriter) isJSONResponse() bool {
	return strings.Contains(w.Header().Get("Content-Type"), "application/json")
}

// Write writes response then backups to ctx
func (w *multiWriter) Write(b []byte) (int, error) {
	var resp logrus.Fields
	if w.isJSONResponse() {
		if err := json.Unmarshal(b, &resp); err != nil {
			return 0, err
		}
		w.ctx.Set("response", resp)
	} else {
		w.ctx.Set("response", b)
	}
	return w.ResponseWriter.Write(b)
}

// WriteString implements StringWriter to writes from gin c.Sting()
func (w *multiWriter) WriteString(b string) (int, error) {
	w.ctx.Set("response", b)
	return w.ResponseWriter.Write([]byte(b))
}

func Milliseconds(t time.Duration) float64 {
	m := t.Seconds() * 1000
	f, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", m), 64)
	return f
}

const maxLengthToFilter = 512

func filterBodyTooLong(fields logrus.Fields) {
	for k, v := range fields {
		switch obj := v.(type) {
		case string:
			if len(obj) > maxLengthToFilter {
				fields[k] = "SIZE(" + strconv.Itoa(len(obj)) + ")"
			}
		case logrus.Fields:
			filterBodyTooLong(obj)
		case []byte:
			if len(obj) > maxLengthToFilter {
				fields[k] = "SIZE(" + strconv.Itoa(len(obj)) + ")"
			}
		default:
			// interface{}
			// Slice类型,并且每个元素类型必须一致
			s := reflect.ValueOf(v)
			if v != nil && s.Kind() == reflect.Slice {
				for i := 0; i < s.Len(); i++ {
					if mp, ok := s.Index(i).Interface().(logrus.Fields); ok {
						filterBodyTooLong(mp)
					}
					if mp, ok := s.Index(i).Interface().(string); ok {
						if len(mp) > maxLengthToFilter {
							if p, ok := s.Index(i).Addr().Interface().(*string); ok {
								*p = "SIZE(" + strconv.Itoa(len(mp)) + ")"
							}
							if p, ok := s.Index(i).Addr().Interface().(*interface{}); ok {
								*p = "SIZE(" + strconv.Itoa(len(mp)) + ")"
							}
						}
					}
				}
			}
		}
	}
}
