package handler

import (
	"context"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/repositories"
	"github.com/google/uuid"
	"time"
)

type TransactionCommandHandler struct {
	writeRepository repositories.TransactionRepository
}

func (h *TransactionCommandHandler) HandleDeposit(ctx context.Context, command domain.CreateDepositTransaction) error {
	event := domain.DepositedTransaction{
		Id:          uuid.New(),
		Transaction: command.Transaction,
		Time:        time.Now(),
	}

	return h.writeRepository.SaveDeposit(ctx, event)
}

func (h *TransactionCommandHandler) HandleWithdraw(ctx context.Context, command domain.CreateWithdrawTransaction) error {
	event := domain.WithdrawnTransaction{
		Id:          uuid.New(),
		Transaction: command.Transaction,
		Time:        time.Now(),
	}

	return h.writeRepository.SaveWithdraw(ctx, event)
}

func NewTransactionCommandHandler(writeRepository repositories.TransactionRepository) TransactionCommandHandler {
	return TransactionCommandHandler{
		writeRepository: writeRepository,
	}
}
