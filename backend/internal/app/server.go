package app

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api/accounts"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api/transactions"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api/transactions/middleware"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	port       int
}

func (s *Server) Start() error {
	port := fmt.Sprintf(":%d", s.port)

	s.router.GET("/live", api.Liveness)
	s.router.GET("/ready", api.Readiness)
	accountsGroup := s.router.Group("/accounts")
	{
		accountsGroup.POST("/", accounts.Create)
		accountsGroup.GET("/:id/balance", accounts.GetBalance)
	}
	transactionsGroup := s.router.Group("/transactions")
	{
		transactionsGroup.POST("/", middleware.LogDepositsOver(decimal.NewFromInt(10000)), transactions.Make) // TODO add config for amount?
	}

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

func New(options ...Options) *Server {
	s := &Server{
		router: gin.Default(),
		port:   8080,
	}

	for _, o := range options {
		o(s)
	}

	return s
}