package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerError struct {
	error
	Status int
	Data   any
}

func (b *HandlerBuilder) errorMiddleware(c *gin.Context) {
	c.Next()

	lastErr := c.Errors.Last()
	if lastErr == nil {
		return
	}
	var handlerErr *HandlerError
	ok := errors.As(lastErr, &handlerErr)

	if !ok || handlerErr.Status == http.StatusInternalServerError {
		c.String(http.StatusInternalServerError, "Something went wrong")
		return
	}
	// TODO delete data?
	/*if err.Data != nil {
		switch err.Data.(type) {
		case string:
			c.String(err.Status, "%s", err.Data)
		default:
			c.JSON(err.Status, err.Data)
		}
		return
	}*/
	c.Status(handlerErr.Status)
}

func (b *HandlerBuilder) authenticate(c *gin.Context) {
	panic("implement me")
}

func (b *HandlerBuilder) authorize(c *gin.Context) {
	panic("implement me")
}
