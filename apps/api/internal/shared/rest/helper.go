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

func BadRequest(c *gin.Context, message string) {
	c.JSON(nethttp.StatusBadRequest, ErrorResponse{
		Error: Error{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(nethttp.StatusNotFound, ErrorResponse{
		Error: Error{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

func Accepted[T any](c *gin.Context, data T) {
	c.JSON(nethttp.StatusAccepted, Response[T]{
		Data: data,
	})
}
