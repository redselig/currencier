package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/redselig/currencier/internal/data/controllers"
	"github.com/redselig/currencier/internal/data/logger/zerologger"
	"github.com/redselig/currencier/internal/data/repository/db"
	"github.com/redselig/currencier/internal/domain/usecase"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(cfg *Config, debug bool) (err error) {

	wr := os.Stdout
	if !debug {
		wr, err = os.OpenFile(cfg.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Wrapf(err, "can't create/open log file")
		}
	}
	logger := zerologger.NewLogger(wr, debug)
	client := controllers.NewHTTPClient(cfg.Update.Source, 30)
	repo, err := db.NewPGSRepo(cfg.DB.Dialect, cfg.DB.DSN)
	if err != nil {
		return errors.Wrap(err, "cant't initialize repository")
	}
	currensier := usecase.NewCurrencierInteractor(client, repo)
	server := controllers.NewHttpServer(net.JoinHostPort("0.0.0.0", cfg.API.HTTPPort), logger, currensier)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Serve(); err != nil {
			logger.Log(ctx, errors.Wrapf(err, "can't start http server"))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.update(ctx, currensier, cfg.Update.Time, logger); err != nil {
			logger.Log(context.Background(), errors.Wrapf(err, "can't start updating currencies"))
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
	server.StopServe()
	wg.Wait()
	return nil
}

func (a *App) update(ctx context.Context, currencier usecase.Currencier, timeout string, logger usecase.Logger) error {
	defer logger.Log(ctx, "stop update currencies in repo")

	d, err := time.ParseDuration(timeout)
	if err != nil {
		return errors.Wrap(err, "cant't parse update timeout")
	}
	ticker := time.NewTicker(d)

	if err := currencier.UpdateCurrencies(ctx); err != nil {
		logger.Log(ctx, errors.Wrapf(err, "can't update currencies in repo"))
	}

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		default:
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := currencier.UpdateCurrencies(ctx); err != nil {
				logger.Log(ctx, errors.Wrapf(err, "can't update currencies in repo"))
			}
		}
	}
}
