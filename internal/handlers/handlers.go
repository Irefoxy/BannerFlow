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

func (b *HandlerBuilder) handleListBanners(c *gin.Context) {
	banners, err := b.listBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.JSON(http.StatusOK, converters.BannersExtToInnerResponses(banners))
}

func (b *HandlerBuilder) handleUserGetBanner(c *gin.Context) {
	content, err := b.userGetBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.JSON(http.StatusOK, content.Content)
}

func (b *HandlerBuilder) handleCreateBanner(c *gin.Context) {
	id, err := b.createBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.JSON(http.StatusCreated, converters.ConstructGet201Response(id))
}

func (b *HandlerBuilder) handleUpdateBanner(c *gin.Context) {
	err := b.updateBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (b *HandlerBuilder) handleDeleteBanner(c *gin.Context) {
	err := b.deleteBanner(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (b *HandlerBuilder) updateBanner(c *gin.Context) error {
	id, err := getAndValidateIntParam(c, idName)
	if err != nil {
		return err
	}
	req, err := readRequest[openapi.BannerGetRequest](c)
	if err != nil {
		return err
	}
	err = validateUpdateBannerRequest(req)
	if err != nil {
		return err
	}
	return b.srv.UpdateBanner(c.Request.Context(), id, converters.GetRequestToUpdateBanner(req))
}

func (b *HandlerBuilder) listBanner(c *gin.Context) ([]models.BannerExt, error) {
	tag, err := getAndValidateIntParam(c, tagName)
	if err != nil {
		tag = models.ZeroValue
	}
	feature, err := getAndValidateIntParam(c, featureName)
	if err != nil {
		feature = models.ZeroValue
	}
	limit, err := getAndValidateIntParam(c, limitName)
	if err != nil {
		limit = models.ZeroValue
	}
	offset, err := getAndValidateIntParam(c, offsetName)
	if err != nil {
		offset = models.ZeroValue
	}
	return b.srv.ListBanners(c.Request.Context(), converters.ConstructBannerListOptions(limit, offset, tag, feature))
}

func (b *HandlerBuilder) createBanner(c *gin.Context) (int, error) {
	req, err := readRequest[openapi.BannerGetRequest](c)
	if err != nil {
		return 0, err
	}
	err = validateCreateBannerRequest(req)
	if err != nil {
		return 0, err
	}
	return b.srv.CreateBanner(c.Request.Context(), &converters.GetRequestToUpdateBanner(req).Banner)
}

func (b *HandlerBuilder) deleteBanner(c *gin.Context) error {
	id, err := getAndValidateIntParam(c, idName)
	if err != nil {
		return err
	}
	return b.srv.DeleteBanner(c.Request.Context(), id)
}

func (b *HandlerBuilder) userGetBanner(c *gin.Context) (*models.UserBanner, error) {
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

func getAndValidateIntParam(c *gin.Context, key string) (int, error) {
	id, ok := c.Params.Get(key)
	if !ok {
		return 0, fmt.Errorf("%w: %s is required", e.ErrorInParam, key)
	}
	iid, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%w: %s should be an int integer", e.ErrorInParam, key)
	}
	return int(iid), nil
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
	case req.TagIds == nil:
		return fmt.Errorf("%w: missing tag ids", e.ErrorInRequestBody)
	case req.FeatureId == nil:
		return fmt.Errorf("%w: missing feature ids", e.ErrorInRequestBody)
	case req.Content == nil:
		return fmt.Errorf("%w: missing content", e.ErrorInRequestBody)
	case req.IsActive == nil:
		return fmt.Errorf("%w: missing is_active", e.ErrorInRequestBody)
	}
	return nil
}
func validateUpdateBannerRequest(req openapi.BannerGetRequest) error {
	if req.TagIds == nil && req.FeatureId == nil && req.Content == nil && req.IsActive == nil {
		return fmt.Errorf("%w: all fields are empty", e.ErrorInRequestBody)
	}
	return nil
}
