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
	FeatureId *int `form:"feature_id" binding:"omitempty,gte=0"`
	TagId     *int `form:"tag_id" binding:"omitempty,gte=0"`
	Limit     *int `form:"limit" binding:"omitempty,gte=1"`
	Offset    *int `form:"offset" binding:"omitempty,gte=0"`
}

type IdParams struct {
	Id int `uri:"id" binding:"required,gt=0"`
}

type BannerRequest struct {
	Content   *map[string]interface{} `json:"content,omitempty" binding:"required"`
	FeatureId *int                    `json:"feature_id,omitempty" binding:"required,gte=0"`
	IsActive  *bool                   `json:"is_active,omitempty" binding:"required"`
	TagIds    *[]int                  `json:"tag_ids,omitempty" binding:"required,gte=1,dive,gte=0"`
}

type BannerUpdateRequest struct {
	Content   *map[string]interface{} `json:"content,omitempty"`
	FeatureId *int                    `json:"feature_id,omitempty" binding:"omitempty,gte=0"`
	IsActive  *bool                   `json:"is_active,omitempty"`
	TagIds    *[]int                  `json:"tag_ids,omitempty" binding:"omitempty,gte=1,dive,gte=0"`
}

type BannerResponse struct {
	BannerId  *int                    `json:"banner_id"`
	Content   *map[string]interface{} `json:"content"`
	CreatedAt *time.Time              `json:"created_at"`
	FeatureId *int                    `json:"feature_id"`
	IsActive  *bool                   `json:"is_active"`
	TagIds    *[]int                  `json:"tag_ids"`
	UpdatedAt *time.Time              `json:"updated_at"`
}

type BannerErrorResponse struct {
	Error *string `json:"error"`
}

type BannerIdResponse struct {
	BannerId *int `json:"banner_id"`
}
