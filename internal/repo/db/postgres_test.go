package db

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/domain/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

var testErr = errors.New("test")

type PostgresTest struct {
	suite.Suite
	pool     pgxmock.PgxPoolIface
	postgres *PostgresDatabase
}

func (s *PostgresTest) SetupTest() {
	var err error
	s.pool, err = pgxmock.NewPool()
	s.Require().NoError(err)
	s.postgres = New(s.pool)
}

func (s *PostgresTest) TestStop() {
	s.pool.ExpectClose()
	s.postgres.Stop()
}

// TestAdd00 tests case ping fails
func (s *PostgresTest) TestAdd00() {
	s.pool.ExpectPing().WillReturnError(testErr)

	gotId, err := s.postgres.Add(context.Background(), nil)
	s.Assert().Equal(0, gotId)
	s.Assert().ErrorIs(err, e.ErrorFailedToConnect)
}

// TestAdd01 tests case begin fails
func (s *PostgresTest) TestAdd01() {
	s.pool.ExpectPing()
	s.pool.ExpectBegin().WillReturnError(testErr)

	gotId, err := s.postgres.Add(context.Background(), nil)
	s.Assert().Equal(0, gotId)
	s.Assert().ErrorIs(err, testErr)
}

// // TestAdd02 tests case query fails
func (s *PostgresTest) TestAdd02() {
	banner := &models.Banner{
		BaseBanner: models.BaseBanner{
			UserBanner: models.UserBanner{
				Content: map[string]any{"title": "test"},
			},
			FeatureId: 1,
			TagIds:    []int{1, 2, 3},
		},
		IsActive: false,
	}

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(insertBannerQuery)).WithArgs(Attrs(banner.Content), banner.TagIds, banner.FeatureId).WillReturnError(testErr)
	s.pool.ExpectRollback()

	gotId, err := s.postgres.Add(context.Background(), banner)
	s.Assert().Equal(0, gotId)
	s.Assert().ErrorIs(err, testErr)
}

// TestAdd03 tests case query conflict
func (s *PostgresTest) TestAdd03() {
	banner := &models.Banner{
		BaseBanner: models.BaseBanner{
			UserBanner: models.UserBanner{
				Content: map[string]any{"title": "test"},
			},
			FeatureId: 1,
			TagIds:    []int{1, 2, 3},
		},
		IsActive: false,
	}
	conflictErr := &pgconn.PgError{
		Message: "conflict",
		Code:    pgErrUniqueViolation,
	}

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(insertBannerQuery)).WithArgs(Attrs(banner.Content), banner.TagIds, banner.FeatureId).WillReturnError(conflictErr)
	s.pool.ExpectRollback()

	gotId, err := s.postgres.Add(context.Background(), banner)
	s.Assert().Equal(0, gotId)
	s.Assert().ErrorIs(err, e.ErrorConflict)
}

// TestAdd04 tests case OK without IsActive false
func (s *PostgresTest) TestAdd04() {
	banner := &models.Banner{
		BaseBanner: models.BaseBanner{
			UserBanner: models.UserBanner{
				Content: map[string]any{"title": "test"},
			},
			FeatureId: 1,
			TagIds:    []int{1, 2, 3},
		},
		IsActive: true,
	}
	const id = 1
	responseRows := pgxmock.NewRows([]string{"id"})
	responseRows.AddRow(id)

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(insertBannerQuery)).WithArgs(Attrs(banner.Content), banner.TagIds, banner.FeatureId).WillReturnRows(responseRows)
	s.pool.ExpectCommit()
	s.pool.ExpectRollback()

	gotId, err := s.postgres.Add(context.Background(), banner)
	s.Assert().Equal(1, gotId)
	s.Assert().NoError(err)
}

// TestAdd05 tests case with IsActive = true and exev fails
func (s *PostgresTest) TestAdd05() {
	banner := &models.Banner{
		BaseBanner: models.BaseBanner{
			UserBanner: models.UserBanner{
				Content: map[string]any{"title": "test"},
			},
			FeatureId: 1,
			TagIds:    []int{1, 2, 3},
		},
		IsActive: false,
	}
	const id = 1
	responseRows := pgxmock.NewRows([]string{"id"})
	responseRows.AddRow(id)

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(insertBannerQuery)).WithArgs(Attrs(banner.Content), banner.TagIds, banner.FeatureId).WillReturnRows(responseRows)
	s.pool.ExpectExec(regexp.QuoteMeta(insertDeactivatedBannerQuery)).WithArgs(id).WillReturnError(testErr)
	s.pool.ExpectRollback()

	gotId, err := s.postgres.Add(context.Background(), banner)
	s.Assert().Equal(0, gotId)
	s.Assert().ErrorIs(err, testErr)
}

// TestAdd05 tests case everything OK
func (s *PostgresTest) TestAdd06() {
	banner := &models.Banner{
		BaseBanner: models.BaseBanner{
			UserBanner: models.UserBanner{
				Content: map[string]any{"title": "test"},
			},
			FeatureId: 1,
			TagIds:    []int{1, 2, 3},
		},
		IsActive: false,
	}
	const id = 1
	responseRows := pgxmock.NewRows([]string{"id"})
	responseRows.AddRow(id)

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(insertBannerQuery)).WithArgs(Attrs(banner.Content), banner.TagIds, banner.FeatureId).WillReturnRows(responseRows)
	s.pool.ExpectExec(regexp.QuoteMeta(insertDeactivatedBannerQuery)).WithArgs(id).WillReturnResult(pgxmock.NewResult("insert", 1))
	s.pool.ExpectCommit()
	s.pool.ExpectRollback()

	gotId, err := s.postgres.Add(context.Background(), banner)
	s.Assert().Equal(1, gotId)
	s.Assert().NoError(err)
}

// TestGetHistory00 tests case ping failed
func (s *PostgresTest) TestGetHistory00() {
	s.pool.ExpectPing().WillReturnError(testErr)

	banners, err := s.postgres.GetHistoryForId(context.Background(), 1)
	s.Assert().ErrorIs(err, e.ErrorFailedToConnect)
	s.Assert().Nil(banners)
}

// TestGetHistory01 tests case query returns err
func (s *PostgresTest) TestGetHistory01() {
	const id = 1
	rows := pgxmock.NewRows([]string{"version", "featureid", "tagids", "content"}).RowError(0, testErr)

	s.pool.ExpectPing()
	s.pool.ExpectQuery(regexp.QuoteMeta(selectHistoryQuery)).WithArgs(id).WillReturnRows(rows)
	banners, err := s.postgres.GetHistoryForId(context.Background(), id)
	s.Assert().ErrorIs(err, testErr)
	s.Assert().Nil(banners)
}

// TestGetHistory02 tests case query returns nothing
func (s *PostgresTest) TestGetHistory02() {
	const id = 1
	rows := pgxmock.NewRows([]string{"version", "featureid", "tagids", "content"})

	s.pool.ExpectPing()
	s.pool.ExpectQuery(regexp.QuoteMeta(selectHistoryQuery)).WithArgs(id).WillReturnRows(rows)
	banners, err := s.postgres.GetHistoryForId(context.Background(), id)
	s.Assert().NoError(err)
	s.Assert().Nil(banners)
}

// TestGetHistory03 tests case OK
func (s *PostgresTest) TestGetHistory03() {
	const id = 1
	expectedBanners := []models.HistoryBanner{
		{
			BaseBanner: models.BaseBanner{
				UserBanner: models.UserBanner{
					Content: map[string]any{"title": "test"},
				},
				FeatureId: 1,
				TagIds:    []int{1, 2, 3},
			},
			Version: 12,
		},
		{
			BaseBanner: models.BaseBanner{
				UserBanner: models.UserBanner{
					Content: map[string]any{"title": "test2"},
				},
				FeatureId: 2,
				TagIds:    []int{1, 2, 3},
			},
			Version: 10,
		},
	}
	rows := pgxmock.NewRows([]string{"version", "featureid", "tagids", "content"})
	for _, banner := range expectedBanners {
		b, err := json.Marshal(banner.Content)
		s.Require().NoError(err)
		rows.AddRow(banner.Version, banner.FeatureId, banner.TagIds, b)
	}

	s.pool.ExpectPing()
	s.pool.ExpectQuery(regexp.QuoteMeta(selectHistoryQuery)).WithArgs(id).WillReturnRows(rows)
	banners, err := s.postgres.GetHistoryForId(context.Background(), id)
	s.Assert().NoError(err)
	s.Assert().ElementsMatch(expectedBanners, banners)
}

// TestSelectBannerVersion00 tests case ping fails
func (s *PostgresTest) TestSelectBannerVersion00() {
	const (
		id      = 0
		version = 12
	)

	s.pool.ExpectPing().WillReturnError(testErr)

	err := s.postgres.SelectBannerVersion(context.Background(), id, version)
	s.Assert().ErrorIs(err, e.ErrorFailedToConnect)
}

// TestSelectBannerVersion01 tests case exec fails
func (s *PostgresTest) TestSelectBannerVersion01() {
	const (
		id      = 0
		version = 12
	)

	s.pool.ExpectPing()
	s.pool.ExpectExec(regexp.QuoteMeta(callSelectVersionProcedure)).WithArgs(id, version).WillReturnError(testErr)

	err := s.postgres.SelectBannerVersion(context.Background(), id, version)
	s.Assert().ErrorIs(err, testErr)
}

// TestSelectBannerVersion02 tests case OK
func (s *PostgresTest) TestSelectBannerVersion02() {
	const (
		id      = 0
		version = 12
	)

	s.pool.ExpectPing()
	s.pool.ExpectExec(regexp.QuoteMeta(callSelectVersionProcedure)).WithArgs(id, version).WillReturnResult(pgxmock.NewResult("exec", 1))

	err := s.postgres.SelectBannerVersion(context.Background(), id, version)
	s.Assert().NoError(err)
}

// TestDelete00 tests case ping fails
func (s *PostgresTest) TestDelete00() {
	ids := []int{1, 2, 3}

	s.pool.ExpectPing().WillReturnError(testErr)

	err := s.postgres.DeleteByIds(context.Background(), ids...)
	s.Assert().ErrorIs(err, e.ErrorFailedToConnect)
}

// TestDelete01 tests case exec fails
func (s *PostgresTest) TestDelete01() {
	ids := []int{1, 2, 3}

	s.pool.ExpectPing()
	s.pool.ExpectExec(regexp.QuoteMeta(deleteBannersQuery)).WithArgs(ids).WillReturnError(testErr)

	err := s.postgres.DeleteByIds(context.Background(), ids...)
	s.Assert().ErrorIs(err, testErr)
}

// TestDelete02 tests case nothing deleted
func (s *PostgresTest) TestDelete02() {
	ids := []int{1, 2, 3}

	s.pool.ExpectPing()
	s.pool.ExpectExec(regexp.QuoteMeta(deleteBannersQuery)).WithArgs(ids).WillReturnResult(pgxmock.NewResult("exec", 0))

	err := s.postgres.DeleteByIds(context.Background(), ids...)
	s.Assert().ErrorIs(err, e.ErrorNotFound)
}

// TestDelete03 tests case OK
func (s *PostgresTest) TestDelete03() {
	ids := []int{1, 2, 3}

	s.pool.ExpectPing()
	s.pool.ExpectExec(regexp.QuoteMeta(deleteBannersQuery)).WithArgs(ids).WillReturnResult(pgxmock.NewResult("exec", 3))

	err := s.postgres.DeleteByIds(context.Background(), ids...)
	s.Assert().NoError(err)
}

// TestSelectBannerVersion00 tests case ping fails
func (s *PostgresTest) TestListBanners00() {
	options := &models.BannerListOptions{
		BannerIdentOptions: models.BannerIdentOptions{
			FeatureId: 1,
			TagId:     1,
		},
		Limit:  12,
		Offset: 10,
	}
	s.pool.ExpectPing().WillReturnError(testErr)

	banners, err := s.postgres.List(context.Background(), options)
	s.Assert().ErrorIs(err, e.ErrorFailedToConnect)
	s.Assert().Nil(banners)
}

// TestSelectBannerVersion01 tests case exec fails
func (s *PostgresTest) TestListBanners01() {
	options := &models.BannerListOptions{
		BannerIdentOptions: models.BannerIdentOptions{
			FeatureId: 1,
			TagId:     1,
		},
		Limit:  12,
		Offset: 10,
	}
	query, args := buildListQuery(options)
	rows := pgxmock.NewRows([]string{"id", "content", "created", "updated", "featureId", "tagIds", "isactive"}).RowError(0, testErr)

	s.pool.ExpectPing()
	s.pool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(rows)

	banners, err := s.postgres.List(context.Background(), options)
	s.Assert().ErrorIs(err, testErr)
	s.Assert().Nil(banners)
}

// TestSelectBannerVersion02 tests case OK
func (s *PostgresTest) TestListBanners02() {
	options := &models.BannerListOptions{
		BannerIdentOptions: models.BannerIdentOptions{
			FeatureId: 1,
			TagId:     1,
		},
		Limit:  12,
		Offset: 10,
	}
	expectedBanners := []models.BannerExt{
		{
			BannerId: 1,
			Banner: models.Banner{
				BaseBanner: models.BaseBanner{
					UserBanner: models.UserBanner{
						Content: map[string]any{"content": "test"},
					},
					FeatureId: 1,
					TagIds:    []int{1, 2, 3},
				},
				IsActive: false,
			},
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
		{
			BannerId: 2,
			Banner: models.Banner{
				BaseBanner: models.BaseBanner{
					UserBanner: models.UserBanner{
						Content: map[string]any{"content": "test"},
					},
					FeatureId: 2,
					TagIds:    []int{1, 2, 3},
				},
				IsActive: true,
			},
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}
	query, args := buildListQuery(options)
	rows := pgxmock.NewRows([]string{"id", "content", "created", "updated", "featureId", "tagIds", "isactive"})
	for _, banner := range expectedBanners {
		rows.AddRow(banner.BannerId, banner.Content, banner.CreatedAt, banner.UpdatedAt, banner.FeatureId, banner.TagIds, banner.IsActive)
	}

	s.pool.ExpectPing()
	s.pool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(rows)

	banners, err := s.postgres.List(context.Background(), options)
	s.Assert().NoError(err)
	s.Assert().ElementsMatch(banners, expectedBanners)
}

// TestDeleteByParam00 tests case ping fails
func (s *PostgresTest) TestDeleteByParam00() {
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     models.ZeroValue,
	}
	s.pool.ExpectPing().WillReturnError(testErr)

	err := s.postgres.DeleteByFeatureOrTag(context.Background(), options)
	s.Assert().ErrorIs(err, e.ErrorFailedToConnect)
}

// TestDeleteByParam01 tests case Begin fails
func (s *PostgresTest) TestDeleteByParam01() {
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     models.ZeroValue,
	}
	s.pool.ExpectPing()
	s.pool.ExpectBegin().WillReturnError(testErr)

	err := s.postgres.DeleteByFeatureOrTag(context.Background(), options)
	s.Assert().ErrorIs(err, testErr)
}

// TestDeleteByParam02 tests case query fails
func (s *PostgresTest) TestDeleteByParam02() {
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     models.ZeroValue,
	}
	query, args := buildDeleteQuery(options)

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(testErr)
	s.pool.ExpectRollback()
	err := s.postgres.DeleteByFeatureOrTag(context.Background(), options)
	s.Assert().ErrorIs(err, testErr)
}

// TestDeleteByParam03 tests case query returns nothing
func (s *PostgresTest) TestDeleteByParam03() {
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     models.ZeroValue,
	}
	query, args := buildDeleteQuery(options)
	row := pgxmock.NewRows([]string{"id"})

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(row)
	s.pool.ExpectRollback()
	err := s.postgres.DeleteByFeatureOrTag(context.Background(), options)
	s.Assert().ErrorIs(err, e.ErrorNotFound)
}

// TestDeleteByParam04 tests case exec fails
func (s *PostgresTest) TestDeleteByParam04() {
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     models.ZeroValue,
	}
	expectedIds := []int{1, 2, 3}

	query, args := buildDeleteQuery(options)
	row := pgxmock.NewRows([]string{"id"})
	row.AddRow(expectedIds)

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(row)
	s.pool.ExpectExec(regexp.QuoteMeta(deleteBannersQuery)).WithArgs(expectedIds).WillReturnError(testErr)
	s.pool.ExpectRollback()

	err := s.postgres.DeleteByFeatureOrTag(context.Background(), options)
	s.Assert().ErrorIs(err, testErr)
}

// TestDeleteByParam05 tests case exec fails
func (s *PostgresTest) TestDeleteByParam05() {
	options := &models.BannerIdentOptions{
		FeatureId: 1,
		TagId:     models.ZeroValue,
	}
	expectedIds := []int{1, 2, 3}

	query, args := buildDeleteQuery(options)
	row := pgxmock.NewRows([]string{"id"})
	row.AddRow(expectedIds)

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(row)
	s.pool.ExpectExec(regexp.QuoteMeta(deleteBannersQuery)).WithArgs(expectedIds).WillReturnResult(pgxmock.NewResult("exec", 3))
	s.pool.ExpectCommit()
	s.pool.ExpectRollback()

	err := s.postgres.DeleteByFeatureOrTag(context.Background(), options)
	s.Assert().NoError(err)
}
func (s *PostgresTest) TearDownTest() {
	s.Assert().NoError(s.pool.ExpectationsWereMet())
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresTest))
}
