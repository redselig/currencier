package mocks

import (
	"context"

	"github.com/redselig/currencier/internal/domain/entity"
)

var _ entity.CurrencyInternalRepository = (*CurrencyInternalRepo)(nil)

type CurrencyInternalRepo struct {
	testCurrency *entity.Currency
}

func (c CurrencyInternalRepo) GetByID(ctx context.Context, id string) (*entity.Currency, error) {
	return c.testCurrency, nil
}

func (c CurrencyInternalRepo) GetPage(ctx context.Context, limit, offset int) ([]*entity.Currency, error) {
	return []*entity.Currency{c.testCurrency}, nil
}

func (c CurrencyInternalRepo) GetLazy(ctx context.Context, limit int, lastID string) ([]*entity.Currency, error) {
	return []*entity.Currency{c.testCurrency}, nil
}

func (c CurrencyInternalRepo) SetAll(ctx context.Context, cs []*entity.Currency) error {
	return nil
}

func NewMockRepo(currency *entity.Currency) *CurrencyInternalRepo {
	return &CurrencyInternalRepo{
		testCurrency: currency,
	}
}
