package handler

import (
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/queries"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type AccountCommandHandler struct {
	writeRepository  repositories.AccountRepository
	transactionQuery queries.TransactionQuery
}

func (h *AccountCommandHandler) HandleCreate(command domain.CreateAccount) error {
	accountCreated := domain.AccountCreated{
		Id:      uuid.New(),
		Account: command.Account,
		Time:    time.Now(),
	}

	return h.writeRepository.Save(accountCreated)
}

func (h *AccountCommandHandler) HandleGetBalance(command domain.GetAccountBalance) (decimal.Decimal, error) {
	// TODO validate account belongs to user
	return h.transactionQuery.GetBalance(command.AccountId), nil
}

func NewAccountCommandHandler(writeRepository repositories.AccountRepository, transactionQuery queries.TransactionQuery) AccountCommandHandler {
	return AccountCommandHandler{
		writeRepository:  writeRepository,
		transactionQuery: transactionQuery,
	}
}
