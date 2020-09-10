package usecase

import (
	"context"

	"github.com/pkg/errors"

	"github.com/redselig/currencier/internal/domain/entity"
)

const (
	ErrLoad   = "can't load currencies"
	ErrGet    = "can't get currency by id"
	ErrGetAll = "can't get currencies"
)

var _ Currencier = (*CurrencierInteractor)(nil)

type CurrencierInteractor struct {
	extRepo entity.CurrencyExternalRepository
	intRepo entity.CurrencyInternalRepository
}

func NewCurrencierInteractor(extRepo entity.CurrencyExternalRepository, intRepo entity.CurrencyInternalRepository) *CurrencierInteractor {
	return &CurrencierInteractor{
		extRepo: extRepo,
		intRepo: intRepo,
	}
}

func (c *CurrencierInteractor) UpdateCurrencies(ctx context.Context) error {
	cs, err := c.extRepo.Load(ctx)
	if err != nil {
		return errors.Wrap(err, ErrLoad)
	}
	err = c.intRepo.SetAll(ctx, cs)
	if err != nil {
		return errors.Wrap(err, ErrLoad)
	}
	return nil
}

func (c *CurrencierInteractor) GetCurrencyBuID(ctx context.Context, id string) (*entity.Currency, error) {
	cr, err := c.intRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, ErrGet)
	}
	return cr, nil
}

func (c *CurrencierInteractor) GetCurrenciesLazy(ctx context.Context, limit int, lastID string) ([]*entity.Currency, error) {
	cs, err := c.intRepo.GetLazy(ctx, limit, lastID)
	if err != nil {
		return nil, errors.Wrap(err, ErrGetAll)
	}
	return cs, nil
}

func (c *CurrencierInteractor) GetCurrenciesPage(ctx context.Context, limit, offset int) ([]*entity.Currency, error) {
	cs, err := c.intRepo.GetPage(ctx, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, ErrGetAll)
	}
	return cs, nil
}
