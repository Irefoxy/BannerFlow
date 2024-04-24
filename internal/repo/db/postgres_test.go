package db

import (
	e "BannerFlow/internal/domain/errors"
	"BannerFlow/internal/domain/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
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

// TestAdd04 tests case with IsActive = true and exev fails
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

func (s *PostgresTest) TearDownTest() {
	s.pool.Close()
	s.Assert().NoError(s.pool.ExpectationsWereMet())
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresTest))
}
