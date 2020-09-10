package usecase

import (
	"context"
	"github.com/redselig/currencier/internal/domain/entity"
)

type Currencier interface {
	UpdateCurrencies(ctx context.Context) error
	GetCurrencyBuID(ctx context.Context, id string) (*entity.Currency, error)
	GetCurrenciesPage(ctx context.Context, limit, offset int) ([]*entity.Currency, error)
	GetCurrenciesLazy(ctx context.Context, limit int, lastID string) ([]*entity.Currency, error)
}
