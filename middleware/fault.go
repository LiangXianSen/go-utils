package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/LiangXianSen/go-utils/errors"
)

// FaultHandler is a middleware for handling errors in each handler
func FaultHandler(c *gin.Context) {
	c.Next()

	if !c.Writer.Written() {
		if errorG, ok := c.Get("error"); ok {
			switch err := errorG.(type) {
			case *errors.Error:
				if err.IsBadRequest() {
					c.JSON(http.StatusBadRequest, err)
				} else {
					c.JSON(http.StatusOK, err)
				}
			default:
				c.JSON(http.StatusInternalServerError, errorG)
			}
		}
	}
}
