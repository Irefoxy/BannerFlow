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

func (s *PostgresTest) TestAdd00() {
	s.pool.ExpectPing().WillReturnError(testErr)

	id, err := s.postgres.Add(context.Background(), nil)
	s.Assert().Equal(0, id)
	s.Assert().ErrorIs(err, e.ErrorFailedToConnect)
}

func (s *PostgresTest) TestAdd01() {
	s.pool.ExpectPing()
	s.pool.ExpectBegin().WillReturnError(testErr)

	id, err := s.postgres.Add(context.Background(), nil)
	s.Assert().Equal(0, id)
	s.Assert().ErrorIs(err, testErr)
}

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
	pgconn.PgError{
		Severity:         "",
		Code:             "",
		Message:          "",
		Detail:           "",
		Hint:             "",
		Position:         0,
		InternalPosition: 0,
		InternalQuery:    "",
		Where:            "",
		SchemaName:       "",
		TableName:        "",
		ColumnName:       "",
		DataTypeName:     "",
		ConstraintName:   "",
		File:             "",
		Line:             0,
		Routine:          "",
	}
	s.pool.ExpectQuery(regexp.QuoteMeta(insertBannerQuery)).WithArgs(Attrs(banner.Content), banner.TagIds, banner.FeatureId).WillReturnError(testErr)

	s.pool.ExpectRollback()

	id, err := s.postgres.Add(context.Background(), banner)
	s.Assert().Equal(0, id)
	s.Assert().ErrorIs(err, testErr)
}

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
	responseRows := pgxmock.NewRows([]string{"id"})
	responseRows.AddRow(1)

	s.pool.ExpectPing()
	s.pool.ExpectBegin()
	s.pool.ExpectQuery(regexp.QuoteMeta(insertBannerQuery)).WithArgs(Attrs(banner.Content), banner.TagIds, banner.FeatureId).WillReturnRows(responseRows)

	s.pool.ExpectRollback()

	id, _ := s.postgres.Add(context.Background(), banner)
	s.Assert().Equal(1, id)
	//s.Assert().ErrorIs(err, testErr)
}

func (s *PostgresTest) TearDownTest() {
	s.pool.Close()
	s.Assert().NoError(s.pool.ExpectationsWereMet())
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresTest))
}
