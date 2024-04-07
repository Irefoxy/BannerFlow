package handlers

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/handlers/converters"
	"BannerFlow/internal/services/models"
	"BannerFlow/pkg/api"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (b *HandlerBuilder) handlerListBanners(c *gin.Context) {
	banners, err := b.handleListBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.JSON(http.StatusOK, converters.BannersExtToInnerResponses(banners))
}

func (b *HandlerBuilder) handlerUserGetBanner(c *gin.Context) {
	content, err := b.handleUserGetBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.JSON(http.StatusOK, content.Content)
}

func (b *HandlerBuilder) handlerCreateBanner(c *gin.Context) {
	id, err := b.handleCreateBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.JSON(http.StatusCreated, converters.ConstructGet201Response(id))
}

func (b *HandlerBuilder) handlerUpdateBanner(c *gin.Context) {
	err := b.handleUpdateBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (b *HandlerBuilder) handlerDeleteBanner(c *gin.Context) {
	err := b.handleDeleteBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (b *HandlerBuilder) handleUpdateBanner(c *gin.Context) error {
	id, err := getAndValidateIntParam(c, idName)
	if err != nil {
		return err
	}
	req, err := readRequest[openapi.BannerGetRequest](c)
	if err != nil {
		return err
	}
	return b.srv.UpdateBanner(c.Request.Context(), id, converters.GetRequestToBanner(req))
}

func (b *HandlerBuilder) handleListBanner(c *gin.Context) ([]models.BannerExt, error) {
	tag, err := getAndValidateIntParam(c, tagName)
	if err != nil {
		tag = zeroValue
	}
	feature, err := getAndValidateIntParam(c, featureName)
	if err != nil {
		feature = zeroValue
	}
	limit, err := getAndValidateIntParam(c, limitName)
	if err != nil {
		limit = zeroValue
	}
	offset, err := getAndValidateIntParam(c, offsetName)
	if err != nil {
		offset = zeroValue
	}
	return b.srv.ListBanners(c.Request.Context(), converters.ConstructBannerListOptions(limit, offset, tag, feature))
}

func (b *HandlerBuilder) handleCreateBanner(c *gin.Context) (int32, error) {
	req, err := readRequest[openapi.BannerGetRequest](c)
	if err != nil {
		return 0, err
	}
	err = validateCreateBannerRequest(req)
	if err != nil {
		return 0, err
	}
	return b.srv.CreateBanner(c.Request.Context(), converters.GetRequestToBanner(req))
}

func (b *HandlerBuilder) handleDeleteBanner(c *gin.Context) error {
	id, err := getAndValidateIntParam(c, idName)
	if err != nil {
		return err
	}
	return b.srv.DeleteBanner(c.Request.Context(), id)
}

func (b *HandlerBuilder) handleUserGetBanner(c *gin.Context) (*models.UserBanner, error) {
	tag, err := getAndValidateIntParam(c, tagName)
	if err != nil {
		return nil, err
	}
	feature, err := getAndValidateIntParam(c, featureName)
	if err != nil {
		return nil, err
	}
	flag, _ := getAndValidateBoolParam(c, lastRevisionFlagName)
	return b.srv.UserGetBanners(c.Request.Context(), converters.ConstructBannerUserOptions(flag, feature, tag))
}

func getAndValidateIntParam(c *gin.Context, key string) (int32, error) {
	id, ok := c.Params.Get(key)
	if !ok {
		return 0, fmt.Errorf("%w: %s is required", e.ErrorInParam, key)
	}
	iid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%w: %s should be an int32 integer", e.ErrorInParam, key)
	}
	return int32(iid), nil
}

func getAndValidateBoolParam(c *gin.Context, key string) (bool, error) {
	id, ok := c.Params.Get(key)
	if !ok {
		return false, fmt.Errorf("%w: %s is required", e.ErrorInParam, key)
	}
	flag, err := strconv.ParseBool(id)
	if err != nil {
		return false, fmt.Errorf("%w: %s should be a bool value", e.ErrorInParam, key)
	}
	return flag, nil
}

func readRequest[T any](c *gin.Context) (T, error) {
	var request T
	if err := c.ShouldBind(&request); err != nil {
		return request, fmt.Errorf("%w: %w", e.ErrorInRequestBody, err)
	}
	return request, nil
}

func collectErrors(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

func validateCreateBannerRequest(req openapi.BannerGetRequest) error {
	switch {
	case !req.HasTagIds():
		return fmt.Errorf("%w: missing tag ids", e.ErrorInRequestBody)
	case !req.HasFeatureId():
		return fmt.Errorf("%w: missing feature ids", e.ErrorInRequestBody)
	case !req.HasContent():
		return fmt.Errorf("%w: missing content", e.ErrorInRequestBody)
	case !req.HasIsActive():
		return fmt.Errorf("%w: missing is_active", e.ErrorInRequestBody)
	}
	return nil
}
