package handlers

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/pkg/api"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (b *HandlerBuilder) errorMiddleware(c *gin.Context) {
	c.Next()

	lastErr := c.Errors.Last()
	if lastErr == nil {
		return
	}
	b.log(c)

	switch {
	case errors.Is(lastErr, e.ErrorNotFound):
		c.Status(http.StatusNotFound)
	case errors.Is(lastErr, e.ErrorNoPermission):
		c.Status(http.StatusForbidden)
	case errors.Is(lastErr, e.ErrorAuthenticationFailed):
		c.Status(http.StatusUnauthorized)
	case errors.Is(lastErr, e.ErrorBadRequest):
		sendJSONError(c, http.StatusBadRequest, lastErr.Error())
	default:
		sendJSONError(c, http.StatusInternalServerError, "something went wrong")
	}
}

func sendJSONError(c *gin.Context, status int, msg string) {
	c.JSON(status, api.BannerErrorResponse{
		Error: &msg,
	})
}

func (b *HandlerBuilder) authenticate(c *gin.Context) {
	err := b.handleAuthentication(c)
	if err != nil {
		collectErrors(c, err)
	}
}

func (b *HandlerBuilder) handleAuthentication(c *gin.Context) error {
	token := &api.TokenParam{}
	err := c.ShouldBindHeader(&token)
	if err != nil {
		return e.ErrorNoToken
	}
	err = b.authenticator.Authenticate(token.Token)
	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrorAuthenticationFailed, err)
	}
	return nil
}

func (b *HandlerBuilder) authorize(c *gin.Context) {
	token := c.GetHeader(tokenName)
	if !b.authorizer.IsAdmin(token) {
		collectErrors(c, e.ErrorNoPermission)
	}
}
