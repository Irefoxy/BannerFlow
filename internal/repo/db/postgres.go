package db

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/services/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
	"time"
)

const (
	bannerTableName     = "banners"
	featureTagTableName = "feature_tag"
	idName              = "id"
	bannerIdName        = "bannerId"
	featureName         = "featureId"
	tagName             = "tagId"
	activeName          = "is_active"
	contentName         = "content"
)

var (
	insertTagFeatureQuery = fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, $3)", featureTagTableName, tagName, featureName, bannerIdName)
	insertBannerQuery     = fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2) RETURNING %s", bannerTableName, contentName, activeName, idName)
	listBannersQuery      = fmt.Sprintf("SELECT * FROM %s", bannerTableName)
)

type IFace interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close()
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func NewPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dsn)
}

type PostgresDatabase struct {
	pool IFace
}

func New(pool IFace) *PostgresDatabase {
	return &PostgresDatabase{
		pool: pool,
	}
}

func (p PostgresDatabase) Add(ctx context.Context, banner *models.Banner) (int, error) {
	if p.pool.Ping(ctx) != nil {
		return 0, e.ErrorFailedToConnect
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var id int
	err = tx.QueryRow(ctx, insertBannerQuery, Attrs(banner.Content), banner.IsActive).Scan(&id)
	if err != nil {
		return 0, err
	}
	for _, tag := range banner.TagId {
		_, err = tx.Exec(ctx, insertTagFeatureQuery, tag, banner.FeatureId, id)
		if err != nil {
			return 0, err
		}
	}
	return id, tx.Commit(ctx)
}

func (p PostgresDatabase) Update(ctx context.Context, id int, banner *models.Banner) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDatabase) List(ctx context.Context, options *models.BannerListOptions) ([]models.BannerExt, error) {
	var banners []models.BannerExt
	query, args := buildQuery(options)
	rows, _ := p.pool.Query(ctx, query, args...)
	var id, featureId int
	var created, updated time.Time
	var isActive bool
	var attr Attrs
	_, err := pgx.ForEachRow(rows, []any{&id, &created, &updated, &featureId, &isActive, &attr}, func() error {
		banners = append(banners, models.BannerExt{
			BannerId: id,
			Banner: models.Banner{
				FeatureId: featureId,
				TagId:     nil,
				IsActive:  isActive,
				UserBanner: models.UserBanner{
					Content: attr,
				},
			},
			UpdatedAt: updated,
			CreatedAt: created,
		})
		return nil
	})
}

func buildQuery(options *models.BannerListOptions) (string, []any) {
	builder := build()
	if options.FeatureId >= 0 || options.TagId >= 0 {
		builder(" WHERE ", nil)
	}
	if options.FeatureId >= 0 {
		builder("feature_id = $", &options.FeatureId)
	}
	if options.TagId >= 0 {
		if options.FeatureId >= 0 {
			builder(" AND ", nil)
		}
		builder("tag_id = $", &options.TagId)
	}
	if options.Limit > 0 {
		builder(" LIMIT $", &options.Limit)
	}
	if options.Offset > 0 {
		builder(" OFFSET $", &options.Offset)
	}
	return builder("", nil)
}

func build() func(str string, arg *int) (string, []any) {
	num := 1
	var args []any
	builder := strings.Builder{}
	builder.WriteString(listBannersQuery)
	return func(str string, arg *int) (string, []any) {
		builder.WriteString(str)
		if arg != nil {
			builder.WriteString(strconv.Itoa(num))
			num++
			args = append(args, *arg)
		}
		return builder.String(), args
	}
}
