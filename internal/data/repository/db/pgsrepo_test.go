package db

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/redselig/currencier/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	testID       = "R01020A"
	testName     = "Азербайджанский манат"
	testRate     = 44.7113
	testCurrency = entity.Currency{
		ID:       testID,
		NumCode:  0,
		CharCode: "",
		Nominal:  1,
		Name:     testName,
		Value:    testRate,
	}
	testLimit  = 3
	testOffset = 1
)

type Suite struct {
	suite.Suite
	repo entity.CurrencyInternalRepository
	mock sqlmock.Sqlmock
	db   *sql.DB
}

func TestRepo(t *testing.T) {
	s := new(Suite)
	suite.Run(t, s)
}

func (s *Suite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.Nil(s.T(), err)
	s.repo = &PGSRepo{s.db}
}

func (s *Suite) TearDownSuite() {
	s.db.Close()
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *Suite) TestPGSRepo_GetByID() {
	ctx := context.TODO()
	s.Run("good test: get currency by id", func() {
		rows := sqlmock.NewRows([]string{"id", "name", "rate"}).
			AddRow(testID, testName, testRate)

		s.mock.ExpectQuery(`select id, name, rate from`).
			WithArgs(testID).
			WillReturnRows(rows)

		c, err := s.repo.GetByID(ctx, testID)
		require.Nil(s.T(), err)
		require.Equal(s.T(), &testCurrency, c)
	})
	s.Run("no rows: get event by id", func() {
		rows := sqlmock.NewRows([]string{"id", "name", "rate"})
		s.mock.ExpectQuery(`select id, name, rate from`).
			WithArgs(testID).
			WillReturnRows(rows)

		c, err := s.repo.GetByID(ctx, testID)

		require.Nil(s.T(), err)
		require.Nil(s.T(), c)
	})
	s.Run("return error: get event by id", func() {
		s.mock.ExpectQuery(`select id, name, rate from`).
			WillReturnError(sql.ErrConnDone)

		c, err := s.repo.GetByID(ctx, testID)

		require.Nil(s.T(), c)
		require.Truef(s.T(), errors.Is(err, sql.ErrConnDone), "GetByID not return cause error")
	})
}

func (s *Suite) TestPGSRepo_GetPage() {
	ctx := context.TODO()
	s.Run("good test: pagination", func() {
		rows := sqlmock.NewRows([]string{"id", "name", "rate"}).
			AddRow(testID, testName, testRate).
			AddRow(testID, testName, testRate).
			AddRow(testID, testName, testRate)


		s.mock.ExpectQuery(`select id, name, rate from`).
			WithArgs(testLimit,testOffset).
			WillReturnRows(rows)

		cs, err := s.repo.GetPage(ctx, testLimit,testOffset)
		require.Nil(s.T(), err)
		require.Equal(s.T(), []*entity.Currency{&testCurrency,&testCurrency,&testCurrency}, cs)
	})
	s.Run("no rows: pagination", func() {
		rows := sqlmock.NewRows([]string{"id", "name", "rate"})
		s.mock.ExpectQuery(`select id, name, rate from`).
			WithArgs(testLimit,testOffset).
			WillReturnRows(rows)

		cs, err := s.repo.GetPage(ctx, testLimit,testOffset)

		require.Nil(s.T(), err)
		require.Nil(s.T(), cs)
	})
	s.Run("return error: pagination", func() {
		s.mock.ExpectQuery(`select id, name, rate from`).
			WillReturnError(sql.ErrConnDone)

		cs, err := s.repo.GetPage(ctx, testLimit,testOffset)

		require.Nil(s.T(), cs)
		require.Truef(s.T(), errors.Is(err, sql.ErrConnDone), "GetByID not return cause error")
	})
}

func (s *Suite) TestPGSRepo_GetLazy() {
	ctx := context.TODO()
	s.Run("good test: lazy load", func() {
		rows := sqlmock.NewRows([]string{"id", "name", "rate"}).
			AddRow(testID, testName, testRate).
			AddRow(testID, testName, testRate).
			AddRow(testID, testName, testRate)


		s.mock.ExpectQuery(`select id, name, rate from`).
			WithArgs(testID,testLimit).
			WillReturnRows(rows)

		cs, err := s.repo.GetLazy(ctx, testLimit,testID)
		require.Nil(s.T(), err)
		require.Equal(s.T(), []*entity.Currency{&testCurrency,&testCurrency,&testCurrency}, cs)
	})
	s.Run("no rows: lazy load", func() {
		rows := sqlmock.NewRows([]string{"id", "name", "rate"})
		s.mock.ExpectQuery(`select id, name, rate from`).
			WithArgs(testID,testLimit).
			WillReturnRows(rows)

		cs, err := s.repo.GetLazy(ctx, testLimit,testID)

		require.Nil(s.T(), err)
		require.Nil(s.T(), cs)
	})
	s.Run("return error: lazy load", func() {
		s.mock.ExpectQuery(`select id, name, rate from`).
			WillReturnError(sql.ErrConnDone)

		cs, err := s.repo.GetLazy(ctx, testLimit,testID)

		require.Nil(s.T(), cs)
		require.Truef(s.T(), errors.Is(err, sql.ErrConnDone), "GetByID not return cause error")
	})
}
func (s *Suite) TestPGSRepo_SetAll() {
	ctx := context.TODO()
	s.Run("good test: save 1 currency to db", func() {
		s.mock.ExpectPrepare(`insert into public.currency`).ExpectExec().
			WithArgs(testID, testName, testRate).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.repo.SetAll(ctx, []*entity.Currency{&testCurrency})
		require.Nil(s.T(), err)
	})
	s.Run("good test: save 0 currency to db", func() {
		s.mock.ExpectPrepare(`insert into public.currency`).ExpectExec().
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := s.repo.SetAll(ctx, []*entity.Currency{})
		require.Nil(s.T(), err)
	})
	s.Run("return error: save currency to db", func() {
		s.mock.ExpectPrepare(`insert into public.currency`).ExpectExec().
			WillReturnError(sql.ErrConnDone)

		err := s.repo.SetAll(ctx, []*entity.Currency{})

		require.Truef(s.T(), errors.Is(err, sql.ErrConnDone), "GetByID not return cause error")
	})
}
