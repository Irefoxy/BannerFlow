package api

// BannerVersionResponse struct to store banner version
type BannerVersionResponse struct {
	Content   *map[string]interface{} `json:"content"`
	TagIds    *[]int                  `json:"tag_ids"`
	FeatureId *int                    `json:"feature_id"`
	Version   *int                    `json:"version"`
}

type DeleteBannerParams struct {
	TagId     *int `form:"tag_id" binding:"omitempty,gte=0"`
	FeatureId *int `form:"feature_id" binding:"omitempty,gte=0"`
}

type SelectBannersParams struct {
	Version int `form:"version" binding:"required,gte=1"`
}
