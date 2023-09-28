package handler

import (
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/repositories"
	"github.com/google/uuid"
	"time"
)

type TransactionCommandHandler struct {
	writeRepository repositories.TransactionRepository
}

func (h *TransactionCommandHandler) HandleDeposit(command domain.CreateDepositTransaction) error {
	event := domain.DepositedTransaction{
		Id:          uuid.New(),
		Transaction: command.Transaction,
		Time:        time.Now(),
	}

	return h.writeRepository.SaveDeposit(event)
}

func (h *TransactionCommandHandler) HandleWithdraw(command domain.CreateWithdrawTransaction) error {
	event := domain.WithdrawnTransaction{
		Id:          uuid.New(),
		Transaction: command.Transaction,
		Time:        time.Now(),
	}

	return h.writeRepository.SaveWithdraw(event)
}

func NewTransactionCommandHandler(writeRepository repositories.TransactionRepository) TransactionCommandHandler {
	return TransactionCommandHandler{
		writeRepository: writeRepository,
	}
}
