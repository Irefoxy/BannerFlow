package db

import (
	"BannerFlow/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"strconv"
	"strings"
)

func buildDeleteQuery(options *models.BannerIdentOptions) (string, []any) {
	builder := build()
	builder(selectIdFromFeatureTagQuery, nil)
	if options.FeatureId > models.ZeroValue {
		builder(" featureId = $", options.FeatureId)
		builder(" GROUP BY featureId", nil)
	} else {
		builder(" tagId = $", options.TagId)
		builder(" GROUP BY tagId", nil)
	}
	return builder("", nil)
}

func buildListQuery(options *models.BannerListOptions) (string, []any) {
	builder := build()
	builder(listBannersQuery, nil)
	featureTagFlag := options.FeatureId > models.ZeroValue || options.TagId > models.ZeroValue
	if featureTagFlag {
		builder(` WHERE EXISTS (SELECT 1 FROM feature_tag ft2 WHERE ft2.bannerId = b.id`, nil)
	}
	if options.FeatureId > models.ZeroValue {
		builder(" AND featureId = $", options.FeatureId)
	}
	if options.TagId > models.ZeroValue {
		builder(" AND tagId = $", options.TagId)
	}
	if featureTagFlag {
		builder(`)`, nil)
	}
	builder(` group by b.id ORDER BY b.id`, nil)
	if options.Limit > models.ZeroValue {
		builder(" LIMIT $", options.Limit)
	}
	if options.Offset > models.ZeroValue {
		builder(" OFFSET $", options.Offset)
	}
	return builder("", nil)
}

func buildUpdateQuery(banner *models.UpdateBanner) (string, []any) {
	builder := build()
	var args []any
	addComma := func() {
		if len(args) > 0 {
			builder(",", nil)
		}
	}
	builder("UPDATE banners SET", nil)
	if banner.Flags&models.FeatureBit > 0 {
		_, args = builder(" featureId=$", banner.FeatureId)
	}
	if banner.Flags&models.TagBit > 0 {
		addComma()
		_, args = builder(" tagIds=$", banner.TagIds)
	}
	if banner.Flags&models.ContentBit > 0 {
		addComma()
		_, args = builder(" content=$", Attrs(banner.Content))
	}
	return builder("", nil)
}

func prepareUpdateBatch(id int, banner *models.UpdateBanner) *pgx.Batch {
	batch := &pgx.Batch{}
	if banner.Flags&models.IsActiveBit > 0 {
		if banner.IsActive {
			batch.Queue(deleteBannerFromDeactivatedQuery, id)
		} else {
			batch.Queue(insertDeactivatedBannerQuery, id)
		}
	}
	if banner.Flags & ^models.IsActiveBit > 0 {
		query, args := buildUpdateQuery(banner)
		batch.Queue(query, args...)
	}
	return batch
}

func build() func(str string, arg any) (string, []any) {
	num := 1
	var args []any
	builder := strings.Builder{}
	return func(str string, arg any) (string, []any) {
		builder.WriteString(str)
		if arg != nil {
			builder.WriteString(strconv.Itoa(num))
			num++
			args = append(args, arg)
		}
		return builder.String(), args
	}
}
