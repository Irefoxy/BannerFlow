package api

import "time"

type AdminParam struct {
	Admin string `uri:"admin"`
}

type TokenParam struct {
	Token string `header:"token" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserBannerParams struct {
	TagId           *int  `form:"tag_id" binding:"required,gte=0"`
	FeatureId       *int  `form:"feature_id" binding:"required,gte=0"`
	UseLastRevision *bool `form:"use_last_revision"`
}

type ListBannerParams struct {
	FeatureId *int `form:"feature_id" binding:"gte=0"`
	TagId     *int `form:"tag_id" binding:"gte=0"`
	Limit     *int `form:"limit" binding:"gte=1"`
	Offset    *int `form:"offset" binding:"gte=0"`
}

type IdParams struct {
	Id int `uri:"id" binding:"required,gt=0"`
}

type BannerRequest struct {
	Content   *map[string]interface{} `json:"content" binding:"required"`
	FeatureId *int                    `json:"feature_id" binding:"required,gte=0"`
	IsActive  *bool                   `json:"is_active" binding:"required"`
	TagIds    *[]int                  `json:"tag_ids" binding:"required,gte=1"`
}

type BannerUpdateRequest struct {
	Content   *map[string]interface{} `json:"content"`
	FeatureId *int                    `json:"feature_id" binding:"gte=0"`
	IsActive  *bool                   `json:"is_active" `
	TagIds    *[]int                  `json:"tag_ids" binding:"gte=1"`
}

type BannerResponse struct {
	BannerId  *int                    `json:"banner_id" binding:"required,gt=0"`
	Content   *map[string]interface{} `json:"content" binding:"required"`
	CreatedAt *time.Time              `json:"created_at" binding:"required"`
	FeatureId *int                    `json:"feature_id" binding:"required,gte=0"`
	IsActive  *bool                   `json:"is_active" binding:"required"`
	TagIds    *[]int                  `json:"tag_ids" binding:"required,gte=1"`
	UpdatedAt *time.Time              `json:"updated_at"`
}

type BannerErrorResponse struct {
	Error *string `json:"error" binding:"required"`
}

type BannerIdResponse struct {
	BannerId *int `json:"banner_id" binding:"required,gt=0"`
}
