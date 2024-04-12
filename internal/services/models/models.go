package models

import "time"

const (
	FeatureBit  = 1
	TagBit      = 1 << 1
	IsActiveBit = 1 << 2
	ContentBit  = 1 << 3
	ZeroValue   = -1
)

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

type UpdateBanner struct {
	Banner
	Flags int
}

type BannerExt struct {
	BannerId int
	Banner
	UpdatedAt time.Time
	CreatedAt time.Time
}
