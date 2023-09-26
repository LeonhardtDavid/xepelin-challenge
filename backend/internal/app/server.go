package app

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api/accounts"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api/middleware"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api/transactions"
	transactMiddleware "github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api/transactions/middleware"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"log/slog"
	"net/http"
)

type Server struct {
	httpServer                *http.Server
	router                    *gin.Engine
	port                      int
	accountCommandHandler     handler.AccountCommandHandler
	transactionCommandHandler handler.TransactionCommandHandler
}

func (s *Server) setupRoutes() {
	// TODO add tracing id for metrics
	s.router.Use(middleware.HandleErrors)

	s.router.GET("/live", api.Liveness)
	s.router.GET("/ready", api.Readiness)
	accountsGroup := s.router.Group("/accounts")
	{
		accountsGroup.Use(middleware.RetrieveCustomer)
		accountsGroup.POST("", accounts.Create(s.accountCommandHandler))
		accountsGroup.GET("/:id/balance", accounts.GetBalance(s.accountCommandHandler))
	}
	transactionsGroup := s.router.Group("/transactions")
	{
		transactionsGroup.Use(middleware.RetrieveCustomer)
		transactionsGroup.POST("",
			transactMiddleware.LogDepositsOver(decimal.NewFromInt(10000)), // TODO add config for amount?
			transactions.Make(s.transactionCommandHandler),
		)
	}
}

func (s *Server) Start(ctx context.Context) error {
	port := fmt.Sprintf(":%d", s.port)

	slog.InfoContext(ctx, "Starting Server at", "port", port)

	s.httpServer = &http.Server{
		Addr:    port,
		Handler: s.router,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

type Options func(server *Server)

func WithPort(port int) Options {
	return func(s *Server) {
		s.port = port
	}
}

func WithAccountCommandHandler(handler handler.AccountCommandHandler) Options {
	return func(s *Server) {
		s.accountCommandHandler = handler
	}
}

func WithTransactionCommandHandler(handler handler.TransactionCommandHandler) Options {
	return func(s *Server) {
		s.transactionCommandHandler = handler
	}
}

func New(options ...Options) *Server {
	s := &Server{
		router: gin.Default(),
		port:   8080,
	}

	for _, o := range options {
		o(s)
	}

	s.setupRoutes()

	return s
}
