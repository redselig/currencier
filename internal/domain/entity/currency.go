package entity

import "context"

type Currency struct {
	ID       string
	NumCode  int
	CharCode string
	Nominal  int
	Name     string
	Value    float64
}

type CurrencyInternalRepository interface {
	GetByID(ctx context.Context,id string) (*Currency, error)
	GetPage(ctx context.Context,limit,offset int) ([]*Currency, error)
	GetLazy(ctx context.Context,limit int,lastID string) ([]*Currency, error)
	SetAll(ctx context.Context,cs []*Currency) error
}
type CurrencyExternalRepository interface {
	Load(context.Context) ([]*Currency, error)
}
