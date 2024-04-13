package converters

import (
	"BannerFlow/internal/domain/models"
	"BannerFlow/pkg/api"
)

func BannerUpdateRequestToUpdateBanner(req *api.BannerUpdateRequest) *models.UpdateBanner {
	flags := models.ZeroBit
	if req.TagIds != nil {
		flags |= models.TagBit
	}
	if req.FeatureId != nil {
		flags |= models.FeatureBit
	}
	if req.Content != nil {
		flags |= models.ContentBit
	}
	if req.IsActive != nil {
		flags |= models.IsActiveBit
	}
	return &models.UpdateBanner{
		Banner: models.Banner{
			BaseBanner: models.BaseBanner{
				FeatureId: getDefaultValue(req.FeatureId),
				TagIds:    getDefaultValue(req.TagIds),
				UserBanner: models.UserBanner{
					Content: getDefaultValue(req.Content),
				},
			},
			IsActive: getDefaultValue(req.IsActive),
		},
		Flags: flags,
	}
}

func BannerRequestToBanner(req *api.BannerRequest) *models.Banner {
	return &models.Banner{
		BaseBanner: models.BaseBanner{
			FeatureId: *req.FeatureId,
			TagIds:    *req.TagIds,
			UserBanner: models.UserBanner{
				Content: *req.Content,
			},
		},
		IsActive: *req.IsActive,
	}
}

func getDefaultValue[T any](ptr *T) (result T) {
	if ptr == nil {
		return
	}
	return *ptr
}

func ConstructBannerUserOptions(params *api.UserBannerParams) *models.BannerUserOptions {
	return &models.BannerUserOptions{
		UseLastRevision: getDefaultValue(params.UseLastRevision),
		BannerIdentOptions: models.BannerIdentOptions{
			FeatureId: *params.FeatureId,
			TagId:     *params.TagId,
		},
	}
}

func ConstructBannerListOptions(params *api.ListBannerParams) *models.BannerListOptions {
	return &models.BannerListOptions{
		BannerIdentOptions: models.BannerIdentOptions{
			FeatureId: setZeroValueIfEmpty(params.FeatureId),
			TagId:     setZeroValueIfEmpty(params.TagId),
		},
		Limit:  setZeroValueIfEmpty(params.Limit),
		Offset: setZeroValueIfEmpty(params.Offset),
	}
}

func setZeroValueIfEmpty(arg *int) int {
	if arg == nil {
		return models.ZeroValue
	}
	return *arg
}

func BannersExtToInnerResponses(banners []models.BannerExt) []api.BannerResponse {
	var result []api.BannerResponse
	for _, banner := range banners {
		result = append(result, api.BannerResponse{
			BannerId:  &banner.BannerId,
			TagIds:    &banner.TagIds,
			FeatureId: &banner.FeatureId,
			Content:   &banner.Content,
			IsActive:  &banner.IsActive,
			CreatedAt: &banner.CreatedAt,
			UpdatedAt: &banner.UpdatedAt,
		})
	}
	return result
}

func ConstructGet201Response(id int) *api.BannerIdResponse {
	return &api.BannerIdResponse{
		BannerId: &id,
	}
}

func HistoryBannersToVersionResponse(banners []models.HistoryBanner) []api.BannerVersionResponse {
	var result []api.BannerVersionResponse
	for _, banner := range banners {
		result = append(result, api.BannerVersionResponse{
			Content:   &banner.Content,
			TagIds:    &banner.TagIds,
			FeatureId: &banner.FeatureId,
			Version:   &banner.Version,
		})
	}
	return result
}

func ConstructIdentOptions(params *api.DeleteBannerParams) *models.BannerIdentOptions {
	return &models.BannerIdentOptions{
		FeatureId: setZeroValueIfEmpty(params.FeatureId),
		TagId:     setZeroValueIfEmpty(params.TagIds),
	}
}
