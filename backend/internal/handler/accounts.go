package handler

import (
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type AccountCommandHandler struct {
	writeRepository       repositories.AccountWriteRepository
	transactionRepository repositories.TransactionReadRepository
}

func (h *AccountCommandHandler) HandleCreate(command domain.CreateAccount) error {
	accountCreated := domain.AccountCreated{
		Id:      uuid.New(),
		Account: command.Account,
		Time:    time.Now(),
	}

	return h.writeRepository.Save(accountCreated)
}

func (h *AccountCommandHandler) HandleGetBalance(command domain.GetAccountBalance) decimal.Decimal {
	return h.transactionRepository.GetBalance(command.AccountId)
}

func NewAccountCommandHandler(writeRepository repositories.AccountWriteRepository, transactionRepository repositories.TransactionReadRepository) AccountCommandHandler {
	return AccountCommandHandler{
		writeRepository:       writeRepository,
		transactionRepository: transactionRepository,
	}
}
