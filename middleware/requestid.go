package middleware

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// bodyWriter implements gin.ResponseWriter
type bodyWriter struct {
	gin.ResponseWriter
	requestID string
}

func (w *bodyWriter) isJSONResponse() bool {
	return strings.Contains(w.Header().Get("Content-Type"), "application/json")
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	nb := b
	if w.isJSONResponse() {
		var body map[string]interface{}
		if err := json.Unmarshal(b, &body); err == nil {
			if _, ok := body["request_id"]; !ok {
				body["request_id"] = w.requestID
				nb, _ = json.Marshal(body)
			}
		}
	}

	return w.ResponseWriter.Write(nb)
}

// RequestID generates uuid for each http handler
func RequestID(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = GenRequestID()
	}
	c.Request.Header.Set("X-Request-Id", requestID)
	c.Set("request_id", requestID)

	writer := &bodyWriter{
		ResponseWriter: c.Writer,
		requestID:      requestID,
	}

	c.Writer = writer
}

// GenRequestID returns uuid
func GenRequestID() string {
	uid, _ := uuid.NewV4()
	return strings.ToUpper(uid.String())
}
