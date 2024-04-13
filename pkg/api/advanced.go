package api

// BannerVersionResponse struct to store banner version
type BannerVersionResponse struct {
	Content   *map[string]interface{} `json:"content"`
	TagIds    *[]int                  `json:"tag_ids"`
	FeatureId *int                    `json:"feature_id"`
	Version   *int                    `json:"version"`
}

type DeleteBannerParams struct {
	TagIds    *int `query:"tag_ids" binding:"required_without=FeatureId,gte=0"`
	FeatureId *int `query:"feature_id" binding:"required_without=TagIds,gte=0"`
}

type SelectBannersParams struct {
	Version int `query:"version" binding:"required,gte=1"`
}
