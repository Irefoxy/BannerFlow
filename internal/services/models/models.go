package models

import "time"

type BannerIdentOptions struct {
	FeatureId int32
	TagId     int32
}

type BannerUserOptions struct {
	BannerIdentOptions
	UseLastRevision bool
}

type BannerListOptions struct {
	BannerIdentOptions
	Limit  int32
	Offset int32
}

type UserBanner struct {
	Content map[string]any
}

type Banner struct {
	FeatureId int32
	TagId     []int32
	IsActive  bool
	UserBanner
}

type BannerExt struct {
	BannerId int32
	Banner
	UpdatedAt time.Time
	CreatedAt time.Time
}
