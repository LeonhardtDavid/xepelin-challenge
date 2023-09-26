package main

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/app"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/handler"
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

	config, err := LoadConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("Error loading configurations: %v", err))
		return
	}

	transactionRepository := repositories.NewDummyTransactionWriteRepository()

	s := app.New(
		app.WithPort(config.Port),
		app.WithAccountCommandHandler(
			handler.NewAccountCommandHandler(
				repositories.NewDummyAccountWriteRepository(),
				repositories.ToDummyTransactionReadRepository(transactionRepository),
			),
		),
		app.WithTransactionCommandHandler(handler.NewTransactionCommandHandler(transactionRepository)),
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
