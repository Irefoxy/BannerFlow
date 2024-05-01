package db

import (
	"BannerFlow/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name         string
		options      *models.BannerIdentOptions
		expectedSql  string
		expectedArgs []interface{}
	}{
		{
			name: "With feature",
			options: &models.BannerIdentOptions{
				FeatureId: 1,
				TagId:     models.ZeroValue,
			},
			expectedSql:  "SELECT ARRAY_AGG(DISTINCT bannerid) ids FROM feature_tag WHERE featureId = $1 GROUP BY featureId",
			expectedArgs: []any{1},
		},
		{
			name: "With tag",
			options: &models.BannerIdentOptions{
				FeatureId: models.ZeroValue,
				TagId:     1,
			},
			expectedSql:  "SELECT ARRAY_AGG(DISTINCT bannerid) ids FROM feature_tag WHERE tagId = $1 GROUP BY tagId",
			expectedArgs: []any{1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			at := assert.New(t)
			sql, args := buildDeleteQuery(test.options)
			at.Equal(test.expectedSql, sql)
			at.Equal(test.expectedArgs, args)
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name         string
		options      *models.BannerListOptions
		expectedSql  string
		expectedArgs []interface{}
	}{
		{
			name: "Empty options",
			options: &models.BannerListOptions{
				BannerIdentOptions: models.BannerIdentOptions{
					FeatureId: models.ZeroValue,
					TagId:     models.ZeroValue,
				},
				Limit:  models.ZeroValue,
				Offset: models.ZeroValue,
			},
			expectedSql: `SELECT b.id, b.content, b.created, b.updated, b.featureId, b.tagIds,
    NOT EXISTS (SELECT 1 FROM deactivated d WHERE d.bannerId = b.id) is_active
	FROM feature_tag ft 
    JOIN banners b on b.id = ft.bannerId group by b.id ORDER BY b.id`,
			expectedArgs: nil,
		},
		{
			name: "With tag and limit",
			options: &models.BannerListOptions{
				BannerIdentOptions: models.BannerIdentOptions{
					FeatureId: models.ZeroValue,
					TagId:     1,
				},
				Limit:  12,
				Offset: models.ZeroValue,
			},
			expectedSql: `SELECT b.id, b.content, b.created, b.updated, b.featureId, b.tagIds,
    NOT EXISTS (SELECT 1 FROM deactivated d WHERE d.bannerId = b.id) is_active
	FROM feature_tag ft 
    JOIN banners b on b.id = ft.bannerId WHERE EXISTS (SELECT 1 FROM feature_tag ft2 WHERE ft2.bannerId = b.id AND tagId = $1) group by b.id ORDER BY b.id LIMIT $2`,
			expectedArgs: []any{1, 12},
		},
		{
			name: "With feature and offset",
			options: &models.BannerListOptions{
				BannerIdentOptions: models.BannerIdentOptions{
					FeatureId: 1,
					TagId:     models.ZeroValue,
				},
				Limit:  models.ZeroValue,
				Offset: 12,
			},
			expectedSql: `SELECT b.id, b.content, b.created, b.updated, b.featureId, b.tagIds,
    NOT EXISTS (SELECT 1 FROM deactivated d WHERE d.bannerId = b.id) is_active
	FROM feature_tag ft 
    JOIN banners b on b.id = ft.bannerId WHERE EXISTS (SELECT 1 FROM feature_tag ft2 WHERE ft2.bannerId = b.id AND featureId = $1) group by b.id ORDER BY b.id OFFSET $2`,
			expectedArgs: []any{1, 12},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			at := assert.New(t)
			sql, args := buildListQuery(test.options)
			at.Equal(test.expectedSql, sql)
			at.Equal(test.expectedArgs, args)
		})
	}
}

func TestUpdate(t *testing.T) {
	const id = 123
	tests := []struct {
		name         string
		banner       *models.UpdateBanner
		expectedSql  string
		expectedArgs []interface{}
	}{
		{
			name: "With feature",
			banner: &models.UpdateBanner{
				Banner: models.Banner{
					BaseBanner: models.BaseBanner{
						UserBanner: models.UserBanner{
							Content: map[string]any{"test": "test"},
						},
						FeatureId: 1,
						TagIds:    []int{1, 2, 3},
					},
				},
				Flags: models.FeatureBit,
			},
			expectedSql:  "UPDATE banners SET featureId=$1 WHERE id = $2",
			expectedArgs: []any{1, id},
		},
		{
			name: "With feature and tags",
			banner: &models.UpdateBanner{
				Banner: models.Banner{
					BaseBanner: models.BaseBanner{
						UserBanner: models.UserBanner{
							Content: map[string]any{"test": "test"},
						},
						FeatureId: 1,
						TagIds:    []int{1, 2, 3},
					},
				},
				Flags: models.FeatureBit | models.TagBit,
			},
			expectedSql:  "UPDATE banners SET featureId=$1, tagIds=$2 WHERE id = $3",
			expectedArgs: []any{1, []int{1, 2, 3}, id},
		},
		{
			name: "With feature, tags and content",
			banner: &models.UpdateBanner{
				Banner: models.Banner{
					BaseBanner: models.BaseBanner{
						UserBanner: models.UserBanner{
							Content: map[string]any{"test": "test"},
						},
						FeatureId: 1,
						TagIds:    []int{1, 2, 3},
					},
				},
				Flags: models.FeatureBit | models.TagBit | models.ContentBit,
			},
			expectedSql:  "UPDATE banners SET featureId=$1, tagIds=$2, content=$3 WHERE id = $4",
			expectedArgs: []any{1, []int{1, 2, 3}, Attrs(map[string]any{"test": "test"}), id},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			at := assert.New(t)
			sql, args := buildUpdateQuery(test.banner, id)
			at.Equal(test.expectedSql, sql)
			at.Equal(test.expectedArgs, args)
		})
	}
}
