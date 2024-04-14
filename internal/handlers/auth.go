package handlers

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/pkg/api"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (b *HandlerBuilder) handleTokenGeneration(c *gin.Context) {
	var token string
	var err error
	param := api.AdminParam{}
	err = c.ShouldBindUri(&param)
	if err == nil && param.Admin == "/admin" {
		token, err = b.generator.GenerateToken(true)
	} else {
		token, err = b.generator.GenerateToken(false)
	}
	if err != nil {
		collectErrors(c, fmt.Errorf("%w: error generating token: %w", e.ErrorInternal, err))
		return
	}
	c.JSON(http.StatusOK, api.TokenResponse{Token: token})
}
