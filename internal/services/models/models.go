package models

import "time"

type BannerIdentOptions struct {
	FeatureId int
	TagId     int
}

type BannerUserOptions struct {
	BannerIdentOptions
	UseLastRevision bool
}

type BannerListOptions struct {
	BannerIdentOptions
	Limit  int
	Offset int
}

type UserBanner struct {
	Content map[string]any
}

type Banner struct {
	FeatureId int
	TagId     []int
	IsActive  bool
	UserBanner
}

type BannerExt struct {
	BannerId int
	Banner
	UpdatedAt time.Time
	CreatedAt time.Time
}
