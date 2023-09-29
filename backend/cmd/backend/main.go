package main

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/app"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/config"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/handler"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/queries"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
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

	dbpool, err := pgxpool.New(ctx, conf.DatabaseUrl)
	if err != nil {
		slog.Error(fmt.Sprintf("Error connecting to the database: %v", err))
		return
	}

	s := app.New(
		app.WithPort(conf.Port),
		app.WithAccountCommandHandler(
			handler.NewAccountCommandHandler(
				repositories.NewPostgresAccountRepository(dbpool),
				queries.NewPostgresAccountQuery(dbpool),
				queries.NewDummyTransactionQuery(dbpool),
			),
		),
		app.WithTransactionCommandHandler(
			handler.NewTransactionCommandHandler(repositories.NewPostgresTransactionRepository(dbpool)),
		),
	)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-sigChan
		if err := s.Stop(ctx); err != nil {
			slog.Error(fmt.Sprintf("Error while shutting down server: %v", err))
		}
		dbpool.Close()
		cancel()
	}()

	err = s.Start(ctx)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Failed to start server due to: %v", err))
	}
}
