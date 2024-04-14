package db

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/domain/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
)

const (
	callSelectVersionProcedure       = "CALL choose_banner_from_history($1,$2)"
	selectHistoryQuery               = "SELECT version, featureid, tagids, content FROM banner_history WHERE bannerid = $1 ORDER BY version"
	selectIdFromFeatureTagQuery      = "SELECT ARRAY_AGG(DISTINCT bannerid) ids FROM feature_tag WHERE"
	deleteBannersQuery               = "DELETE FROM banners WHERE id = ANY($1)"
	deleteBannerFromDeactivatedQuery = "DELETE FROM deactivated WHERE bannerid = $1"
	insertDeactivatedBannerQuery     = "INSERT INTO deactivated (bannerid) VALUES ($1) ON CONFLICT DO NOTHING"
	insertBannerQuery                = "INSERT INTO banners (content, tagIds, featureId) VALUES ($1, $2, $3) RETURNING id"
	listBannersQuery                 = `SELECT b.id, b.content, b.created, b.updated, b.featureId, b.tagIds,
    EXISTS (SELECT 1 FROM deactivated d WHERE d.bannerId = b.id) is_active
	FROM feature_tag ft 
    JOIN banners b on b.id = ft.bannerId`
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
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

type PostgresDatabase struct {
	pool IFace
}

func New(pool IFace) *PostgresDatabase {
	return &PostgresDatabase{
		pool: pool,
	}
}

func (p PostgresDatabase) Stop() {
	p.pool.Close()
}

func (p PostgresDatabase) Add(ctx context.Context, banner *models.Banner) (int, error) {
	if err := p.pool.Ping(ctx); err != nil {
		return 0, e.ErrorFailedToConnect
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var id int
	err = tx.QueryRow(ctx, insertBannerQuery, Attrs(banner.Content), banner.TagIds, banner.FeatureId).Scan(&id)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return 0, e.ErrorConflict
		}
		return 0, err
	}
	if !banner.IsActive {
		_, err = tx.Exec(ctx, insertDeactivatedBannerQuery, id)
		if err != nil {
			return 0, err
		}
	}
	return id, tx.Commit(ctx)
}

func (p PostgresDatabase) Update(ctx context.Context, id int, banner *models.UpdateBanner) error {
	if p.pool.Ping(ctx) != nil {
		return e.ErrorFailedToConnect
	}
	_, err := p.pool.SendBatch(ctx, prepareUpdateBatch(id, banner)).Exec()
	return err
}

func (p PostgresDatabase) GetHistoryForId(ctx context.Context, id int) ([]models.HistoryBanner, error) {
	rows, _ := p.pool.Query(ctx, selectHistoryQuery, id)
	historyBanners, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.HistoryBanner, error) {
		res := models.HistoryBanner{}
		var attr Attrs
		err := row.Scan(&res.Version, &res.FeatureId, &res.TagIds, attr)
		res.Content = attr
		return res, err
	})
	return historyBanners, err
}

func (p PostgresDatabase) SelectBannerVersion(ctx context.Context, id, version int) error {
	_, err := p.pool.Exec(ctx, callSelectVersionProcedure, id, version)
	return err
}

func (p PostgresDatabase) DeleteByIds(ctx context.Context, ids ...int) error {
	if p.pool.Ping(ctx) != nil {
		return e.ErrorFailedToConnect
	}
	tag, err := p.pool.Exec(ctx, deleteBannersQuery, ids)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return e.ErrorNotFound
	}
	return nil
}

func (p PostgresDatabase) DeleteByFeatureOrTag(ctx context.Context, options *models.BannerIdentOptions) error {
	if p.pool.Ping(ctx) != nil {
		return e.ErrorFailedToConnect
	}
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	query, args := buildDeleteQuery(options)
	var ids []int
	err = tx.QueryRow(ctx, query, args...).Scan(&ids)
	if errors.Is(err, pgx.ErrNoRows) {
		return e.ErrorNotFound
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, deleteBannersQuery, ids)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (p PostgresDatabase) List(ctx context.Context, options *models.BannerListOptions) ([]models.BannerExt, error) {
	if p.pool.Ping(ctx) != nil {
		return nil, e.ErrorFailedToConnect
	}
	query, args := buildListQuery(options)
	rows, _ := p.pool.Query(ctx, query, args...)
	banners, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.BannerExt, error) {
		res := models.BannerExt{}
		var attr Attrs
		err := row.Scan(&res.BannerId, &attr, &res.CreatedAt, &res.UpdatedAt, &res.FeatureId, &res.TagIds, &res.IsActive)
		res.Content = attr
		return res, err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return banners, err
}

func buildDeleteQuery(options *models.BannerIdentOptions) (string, []any) {
	builder := build()
	builder(selectIdFromFeatureTagQuery, "")
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
	if options.FeatureId > models.ZeroValue || options.TagId > models.ZeroValue {
		builder(`WHERE EXISTS (
    				SELECT 1 FROM feature_tag ft2
    				WHERE ft2.bannerId = b.id`, nil)
	}
	if options.FeatureId > models.ZeroValue {
		builder(" AND featureId = $", options.FeatureId)
	}
	if options.TagId > models.ZeroValue {
		builder(" AND tagId = $", options.TagId)
	}
	builder(`)
				group by b.id
				ORDER BY b.id`, nil)
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
	if banner.Flags & ^models.FeatureBit > 0 {
		_, args = builder(" featureId=$", banner.FeatureId)
	}
	addComma()
	if banner.Flags & ^models.TagBit > 0 {
		_, args = builder(" tagIds=$", banner.TagIds)
	}
	addComma()
	if banner.Flags & ^models.ContentBit > 0 {
		_, args = builder(" content=$", Attrs(banner.Content))
	}
	return builder("", nil)
}

func prepareUpdateBatch(id int, banner *models.UpdateBanner) *pgx.Batch {
	batch := &pgx.Batch{}
	if banner.Flags&models.IsActiveBit > 0 {
		if banner.IsActive {
			batch.Queue(insertDeactivatedBannerQuery, id)
		} else {
			batch.Queue(deleteBannerFromDeactivatedQuery, id)
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
	builder.WriteString(listBannersQuery)
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
