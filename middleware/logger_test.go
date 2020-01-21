package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/LiangXianSen/go-utils/common"
)

func TestLoggerM(t *testing.T) {
	must := assert.New(t)

	loggerOpts := LoggerOptions{
		Application:  "tester",
		Version:      "v0.0.1",
		EnableOutput: true,
		EnableDebug:  true,
	}

	s := gin.New()
	s.Use(LoggerM(loggerOpts))
	s.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	s.ServeHTTP(w, r)

	must.Equal(http.StatusOK, w.Code)
	must.Equal("OK", w.Body.String())
}

func TestLoggerMWithFormData(t *testing.T) {
	must := assert.New(t)

	loggerOpts := LoggerOptions{
		Application:  "tester",
		Version:      "v0.0.1",
		EnableOutput: true,
		EnableDebug:  true,
	}

	output := common.CaptureStdout(func() {
		s := gin.New()
		s.Use(LoggerM(loggerOpts))
		s.GET("/", func(c *gin.Context) {
			c.String(200, "OK")
		})

		form := url.Values{}
		form.Add("k1", "v1")
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(form.Encode()))
		s.ServeHTTP(w, r)
		must.Equal(http.StatusOK, w.Code)
		must.Equal("OK", w.Body.String())
	})

	fmt.Println(string(output))
	var logs map[string]interface{}
	err := json.Unmarshal(output, &logs)
	must.Nil(err)
	must.Equal("k1=v1", logs["request_body"].(string))
}

func TestLoggerMWithQueryParams(t *testing.T) {
	must := assert.New(t)

	loggerOpts := LoggerOptions{
		Application:  "tester",
		Version:      "v0.0.1",
		EnableOutput: true,
		EnableDebug:  true,
	}

	output := common.CaptureStdout(func() {
		s := gin.New()
		s.Use(LoggerM(loggerOpts))
		s.GET("/", func(c *gin.Context) {
			c.String(200, "OK")
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?k1=v1", nil)
		s.ServeHTTP(w, r)
		must.Equal(http.StatusOK, w.Code)
		must.Equal("OK", w.Body.String())
	})

	fmt.Println(string(output))
	var logs map[string]interface{}
	err := json.Unmarshal(output, &logs)
	must.Nil(err)
	must.NotNil(logs["params"])
}

func TestLoggerRecords(t *testing.T) {
	must := assert.New(t)

	loggerOpts := LoggerOptions{
		Application:  "tester",
		Version:      "v0.0.1",
		EnableOutput: true,
		EnableDebug:  true,
	}

	output := common.CaptureStdout(func() {
		s := gin.New()
		s.Use(LoggerM(loggerOpts))
		s.GET("/", func(c *gin.Context) {
			logger := c.MustGet("logger").(*Logger)
			logger.Records("k1", "v1")
			c.String(200, "OK")
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		s.ServeHTTP(w, r)
		must.Equal(http.StatusOK, w.Code)
		must.Equal("OK", w.Body.String())
	})

	fmt.Println(string(output))
	var logs map[string]interface{}
	err := json.Unmarshal(output, &logs)
	must.Nil(err)
	must.NotNil(logs["pipeline"])
}
