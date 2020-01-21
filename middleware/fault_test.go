package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	PErr "github.com/LiangXianSen/go-utils/errors"
)

func TestErrorHandlingWithBadRequest(t *testing.T) {
	must := assert.New(t)

	s := gin.New()
	s.Use(FaultHandler)
	s.GET("/", func(c *gin.Context) {
		err := errors.New("empty params")
		c.Set("error", PErr.InvalidParamError.Wrap(err))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	s.ServeHTTP(w, r)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	must.Nil(err)

	must.Equal(float64(1002), resp["code"].(float64))
	must.Equal(http.StatusBadRequest, w.Code)
}

func TestErrorHandlingWithInternalError(t *testing.T) {
	must := assert.New(t)

	s := gin.New()
	s.Use(FaultHandler)
	s.GET("/", func(c *gin.Context) {
		err := errors.New("Cannot continue executing code")
		c.Set("error", PErr.InternalError.Wrap(err))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	s.ServeHTTP(w, r)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	must.Nil(err)

	must.Equal(float64(1007), resp["code"].(float64))
	must.Equal(http.StatusOK, w.Code)
}
