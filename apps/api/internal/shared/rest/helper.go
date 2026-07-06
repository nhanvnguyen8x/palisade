package rest

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
)

func OK[T any](c *gin.Context, data T) {
	c.JSON(nethttp.StatusOK, Response[T]{
		Data: data,
	})
}

func InternalServerError(c *gin.Context) {
	c.JSON(nethttp.StatusInternalServerError, ErrorResponse{
		Error: Error{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Unexpected error occurred.",
		},
	})
}
