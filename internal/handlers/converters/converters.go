package converters

import (
	"BannerFlow/internal/services/models"
	"BannerFlow/pkg/api"
)

func GetRequestToUpdateBanner(req openapi.BannerGetRequest) *models.UpdateBanner {
	flags := 0
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
			FeatureId: getDefaultValue(req.FeatureId),
			TagId:     getDefaultValue(req.TagIds),
			IsActive:  getDefaultValue(req.IsActive),
			UserBanner: models.UserBanner{
				Content: getDefaultValue(req.Content),
			},
		},
		Flags: flags,
	}
}

func getDefaultValue[T any](ptr *T) (result T) {
	if ptr == nil {
		return
	}
	return *ptr
}

func ConstructBannerUserOptions(flag bool, feature, tag int) *models.BannerUserOptions {
	return &models.BannerUserOptions{
		UseLastRevision: flag,
		BannerIdentOptions: models.BannerIdentOptions{
			FeatureId: feature,
			TagId:     tag,
		},
	}
}

func ConstructBannerListOptions(limit, offset, feature, tag int) *models.BannerListOptions {
	return &models.BannerListOptions{
		BannerIdentOptions: models.BannerIdentOptions{
			FeatureId: feature,
			TagId:     tag,
		},
		Limit:  limit,
		Offset: offset,
	}
}

func BannersExtToInnerResponses(banners []models.BannerExt) []openapi.BannerGet200ResponseInner {
	var result []openapi.BannerGet200ResponseInner
	for _, banner := range banners {
		result = append(result, openapi.BannerGet200ResponseInner{
			BannerId:  &banner.BannerId,
			TagIds:    &banner.TagId,
			FeatureId: &banner.FeatureId,
			Content:   &banner.Content,
			IsActive:  &banner.IsActive,
			CreatedAt: &banner.CreatedAt,
			UpdatedAt: &banner.UpdatedAt,
		})
	}
	return result
}

func ConstructGet201Response(id int) *openapi.BannerGet201Response {
	return &openapi.BannerGet201Response{
		BannerId: &id,
	}
}
