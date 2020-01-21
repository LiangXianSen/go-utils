package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	must := assert.New(t)

	s := gin.New()
	s.Use(RequestID)
	s.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	s.ServeHTTP(w, r)

	fmt.Println(w.Body.String())

	must.Equal(http.StatusOK, w.Code)
}

func TestRequestIDWithEmptyField(t *testing.T) {
	must := assert.New(t)

	s := gin.New()
	s.Use(RequestID)
	s.GET("/", func(c *gin.Context) {
		var resp = struct {
			Status    string `json:"status"`
			RequestID string `json:"request_id,omitempty"`
		}{
			Status: "OK",
		}
		c.JSON(200, resp)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	s.ServeHTTP(w, r)

	fmt.Println(w.Body.String())

	must.Equal(http.StatusOK, w.Code)
}
