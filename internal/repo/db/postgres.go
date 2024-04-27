package db

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/domain/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pgErrUniqueViolation             = "23505"
	callSelectVersionProcedure       = "CALL choose_banner_from_history($1,$2)"
	selectHistoryQuery               = "SELECT version, featureid, tagids, content FROM banner_history WHERE bannerid = $1 ORDER BY version"
	selectIdFromFeatureTagQuery      = "SELECT ARRAY_AGG(DISTINCT bannerid) ids FROM feature_tag WHERE"
	deleteBannersQuery               = "DELETE FROM banners WHERE id = ANY($1)"
	deleteBannerFromDeactivatedQuery = "DELETE FROM deactivated WHERE bannerid = $1"
	insertDeactivatedBannerQuery     = "INSERT INTO deactivated (bannerid) VALUES ($1)"
	insertBannerQuery                = "INSERT INTO banners (content, tagIds, featureId) VALUES ($1, $2, $3) RETURNING id"
	listBannersQuery                 = `SELECT b.id, b.content, b.created, b.updated, b.featureId, b.tagIds,
    NOT EXISTS (SELECT 1 FROM deactivated d WHERE d.bannerId = b.id) is_active
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
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
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
	batch := prepareUpdateBatch(id, banner)
	br := p.pool.SendBatch(ctx, batch)
	for i := 0; i < batch.Len(); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}

func (p PostgresDatabase) GetHistoryForId(ctx context.Context, id int) ([]models.HistoryBanner, error) {
	if p.pool.Ping(ctx) != nil {
		return nil, e.ErrorFailedToConnect
	}
	rows, _ := p.pool.Query(ctx, selectHistoryQuery, id)
	historyBanners, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.HistoryBanner, error) {
		res := models.HistoryBanner{}
		var attr Attrs
		err := row.Scan(&res.Version, &res.FeatureId, &res.TagIds, &attr)
		res.Content = attr
		return res, err
	})
	if err != nil {
		return nil, err
	}
	if len(historyBanners) == 0 {
		return nil, nil
	}
	return historyBanners, nil
}

func (p PostgresDatabase) SelectBannerVersion(ctx context.Context, id, version int) error {
	if p.pool.Ping(ctx) != nil {
		return e.ErrorFailedToConnect
	}
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
		attr := make(Attrs)
		err := row.Scan(&res.BannerId, &attr, &res.CreatedAt, &res.UpdatedAt, &res.FeatureId, &res.TagIds, &res.IsActive)
		res.Content = attr
		return res, err
	})
	if err != nil {
		return nil, err
	}
	if len(banners) == 0 {
		return nil, nil
	}
	return banners, nil
}
