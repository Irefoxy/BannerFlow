package handlers

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/domain/models"
	"BannerFlow/internal/handlers/converters"
	"BannerFlow/pkg/api"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (b *HandlerBuilder) handleDeleteBannerByTagOrFeature(c *gin.Context) {
	err := b.DeleteBannerByTagOrFeature(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.Status(http.StatusAccepted)
}

func (b *HandlerBuilder) handleListBannerHistory(c *gin.Context) {
	banners, err := b.ListBannerHistory(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.JSON(http.StatusOK, converters.HistoryBannersToVersionResponse(banners))
}

func (b *HandlerBuilder) handleSelectBannerVersion(c *gin.Context) {
	err := b.selectBannerVersion(c)
	if err != nil {
		collectErrors(c, err)
		return
	}
	c.Status(http.StatusOK)
}

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

func (b *HandlerBuilder) DeleteBannerByTagOrFeature(c *gin.Context) error {
	params := &api.DeleteBannerParams{}
	err := c.ShouldBindQuery(params)
	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrorInParam, err)
	}
	err = validateDeleteBannerByTagOrFeature(params)
	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrorInParam, err)
	}
	return b.srv.DeleteBannersByTagOrFeature(c.Request.Context(), converters.ConstructIdentOptions(params))
}

func (b *HandlerBuilder) selectBannerVersion(c *gin.Context) error {
	id := &api.IdParams{}
	err := c.ShouldBindUri(id)
	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrorInParam, err)
	}
	version := &api.SelectBannersParams{}
	err = c.ShouldBindQuery(version)
	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrorInParam, err)
	}
	return b.srv.SelectBannerVersion(c.Request.Context(), id.Id, version.Version)
}

func (b *HandlerBuilder) ListBannerHistory(c *gin.Context) ([]models.HistoryBanner, error) {
	id := &api.IdParams{}
	err := c.ShouldBindUri(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", e.ErrorInParam, err)
	}
	return b.srv.ListBannerHistory(c.Request.Context(), id.Id)
}

func (b *HandlerBuilder) updateBanner(c *gin.Context) error {
	id := &api.IdParams{}
	err := c.ShouldBindUri(id)
	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrorInParam, err)
	}
	req, err := readRequest[api.BannerUpdateRequest](c)
	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrorInRequestBody, err)
	}
	updateBanner := converters.BannerUpdateRequestToUpdateBanner(req)
	if updateBanner.Flags == models.ZeroBit {
		return fmt.Errorf("%w: all fields are empty", e.ErrorInRequestBody)
	}
	return b.srv.UpdateBanner(c.Request.Context(), id.Id, updateBanner)
}

func (b *HandlerBuilder) listBanner(c *gin.Context) ([]models.BannerExt, error) {
	params := &api.ListBannerParams{}
	err := c.ShouldBindQuery(params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", e.ErrorInParam, err)
	}
	return b.srv.ListBanners(c.Request.Context(), converters.ConstructBannerListOptions(params))
}

func (b *HandlerBuilder) createBanner(c *gin.Context) (int, error) {
	req, err := readRequest[api.BannerRequest](c)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", e.ErrorInRequestBody, err)
	}
	return b.srv.CreateBanner(c.Request.Context(), converters.BannerRequestToBanner(req))
}

func (b *HandlerBuilder) deleteBanner(c *gin.Context) error {
	id := &api.IdParams{}
	err := c.ShouldBindUri(id)
	if err != nil {
		return fmt.Errorf("%w: %w", e.ErrorInParam, err)
	}
	return b.srv.DeleteBanner(c.Request.Context(), id.Id)
}

func (b *HandlerBuilder) userGetBanner(c *gin.Context) (*models.UserBanner, error) {
	params := &api.UserBannerParams{}
	err := c.ShouldBindQuery(params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", e.ErrorInParam, err)
	}
	return b.srv.UserGetBanners(c.Request.Context(), converters.ConstructBannerUserOptions(params))
}

func readRequest[T any](c *gin.Context) (*T, error) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		return nil, err
	}
	return &request, nil
}

func validateDeleteBannerByTagOrFeature(params *api.DeleteBannerParams) error {
	tagExists := params.TagId != nil
	featureExists := params.FeatureId != nil
	if tagExists != featureExists {
		return nil
	}
	return errors.New("tag id, feature id both exist or not exist")
}

func collectErrors(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}
