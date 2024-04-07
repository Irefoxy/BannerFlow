package converters

import (
	"BannerFlow/internal/services/models"
	"BannerFlow/pkg/api"
)

func GetRequestToBanner(req openapi.BannerGetRequest) *models.Banner {
	return &models.Banner{
		FeatureId: req.GetFeatureId(),
		TagId:     req.GetTagIds(),
		IsActive:  req.GetIsActive(),
		UserBanner: models.UserBanner{
			Content: req.GetContent(),
		},
	}
}

func ConstructBannerUserOptions(flag bool, feature, tag int32) *models.BannerUserOptions {
	return &models.BannerUserOptions{
		UseLastRevision: flag,
		BannerIdentOptions: models.BannerIdentOptions{
			FeatureId: feature,
			TagId:     tag,
		},
	}
}

func ConstructBannerListOptions(limit, offset, feature, tag int32) *models.BannerListOptions {
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
			TagIds:    banner.TagId,
			FeatureId: &banner.FeatureId,
			Content:   banner.Content,
			IsActive:  &banner.IsActive,
			CreatedAt: &banner.CreatedAt,
			UpdatedAt: &banner.UpdatedAt,
		})
	}
	return result
}

func ConstructGet201Response(id int32) *openapi.BannerGet201Response {
	return &openapi.BannerGet201Response{
		BannerId: &id,
	}
}
