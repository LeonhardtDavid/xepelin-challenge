package main

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/app"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/config"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/handler"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/infra"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/queries"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/repositories"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.LoadConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("Error loading configurations: %v", err))
		return
	}

	accountStorage := infra.DummyAccountStorage{}
	transactionStorage := infra.DummyTransactionStorage{}

	s := app.New(
		app.WithPort(conf.Port),
		app.WithAccountCommandHandler(
			handler.NewAccountCommandHandler(
				repositories.NewDummyAccountRepository(&accountStorage),
				queries.NewDummyAccountQuery(&accountStorage),
				queries.NewDummyTransactionQuery(&transactionStorage),
			),
		),
		app.WithTransactionCommandHandler(
			handler.NewTransactionCommandHandler(repositories.NewDummyTransactionRepository(&transactionStorage)),
		),
	)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sigChan
		if err := s.Stop(ctx); err != nil {
			slog.Error(fmt.Sprintf("Error while shutting down server: %v", err))
		}
		cancel()
	}()

	err = s.Start(ctx)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to start server due to: %v", err))
	}
}
