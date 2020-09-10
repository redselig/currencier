package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/redselig/currencier/internal/domain/entity"
)

const (
	ErrAdd = "can't add new rows to table"
	ErrGet = "can't get currencies from db"
)

var _ entity.CurrencyInternalRepository = (*PGSRepo)(nil)

type PGSRepo struct {
	db *sql.DB
}

func NewPGSRepo(driver, dsn string) (*PGSRepo, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create connect to db with dsn %v by driver %v", dsn, driver)
	}
	return &PGSRepo{
		db: db,
	}, nil
}

func (repo *PGSRepo) SetAll(ctx context.Context, cs []*entity.Currency) error {
	sqlStr := "insert into public.currency (id, name, rate) values "
	var vals []interface{}
	for i, row := range cs {
		sqlStr += fmt.Sprintf("($%v,$%v,$%v),", (i*3)+1, (i*3)+2, (i*3)+3)
		vals = append(vals, row.ID, row.Name, row.Value/float64(row.Nominal))
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	sqlStr += " on conflict (id) do UPDATE SET (rate,insert_dt)=(EXCLUDED.rate,now());"
	stmt, err := repo.db.Prepare(sqlStr)
	if err != nil {
		return errors.Wrapf(err, ErrAdd)
	}

	result, err := stmt.ExecContext(ctx, vals...)
	if err != nil {
		return errors.Wrapf(err, ErrAdd)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, ErrAdd)
	}

	if rows != int64(len(cs)) {
		return errors.Wrapf(err, ErrAdd)
	}
	return nil
}

func (repo *PGSRepo) GetByID(ctx context.Context, id string) (*entity.Currency, error) {
	row := repo.db.QueryRowContext(ctx, `select id, name, rate 
												from public.currency where id=$1;`, id)
	if row == nil {
		return nil, nil
	}

	c := entity.Currency{Nominal: 1}
	err := row.Scan(&c.ID, &c.Name, &c.Value)
	if err != nil {
		return nil, SQLError(err, ErrGet)
	}

	return &c, nil
}

func (repo *PGSRepo) GetPage(ctx context.Context, limit, offset int) ([]*entity.Currency, error) {
	rows, err := repo.db.QueryContext(ctx, `select id, name, rate 
												from public.currency order by id limit $1 offset $2;`, limit, offset)
	if err != nil && err != sql.ErrNoRows {
		return nil, SQLError(err, ErrGet)
	}
	defer rows.Close()
	return repo.rowsToCurrencies(rows, ErrGet)
}

func (repo *PGSRepo) GetLazy(ctx context.Context, limit int, lastID string) ([]*entity.Currency, error) {
	rows, err := repo.db.QueryContext(ctx, `select id, name, rate
												from public.currency where id>$1 order by id limit $2;`, lastID, limit)
	if err != nil && err != sql.ErrNoRows {
		return nil, SQLError(err, ErrGet)
	}
	defer rows.Close()
	return repo.rowsToCurrencies(rows, ErrGet)
}

func (repo *PGSRepo) Connect(ctx context.Context, dsn string) (err error) {
	err = repo.db.PingContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to db: %v", dsn)
	}
	return nil
}

func (repo *PGSRepo) Close() error {
	return repo.db.Close()
}

func (repo *PGSRepo) rowsToCurrencies(rows *sql.Rows, errorString string) ([]*entity.Currency, error) {
	var currencies []*entity.Currency
	for rows.Next() {
		c := entity.Currency{Nominal: 1}
		err := rows.Scan(&c.ID, &c.Name, &c.Value)
		if err != nil {
			return nil, SQLError(err, errorString)
		}
		currencies = append(currencies, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, SQLError(err, errorString)
	}
	return currencies, nil
}

func SQLError(err error, message string) error {
	switch err {
	case sql.ErrNoRows:
		return nil
	default:
		return errors.Wrap(err, message)
	}
}
