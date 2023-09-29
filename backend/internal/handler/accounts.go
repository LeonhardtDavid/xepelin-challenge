package handler

import (
	"context"
	"errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/queries"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type AccountCommandHandler struct {
	accountRepository repositories.AccountRepository
	accountQuery      queries.AccountQuery
	transactionQuery  queries.TransactionQuery
}

func (h *AccountCommandHandler) HandleCreate(ctx context.Context, command domain.CreateAccount) error {
	accountCreated := domain.AccountCreated{
		Id:      uuid.New(),
		Account: command.Account,
		Time:    time.Now(),
	}

	return h.accountRepository.Save(ctx, accountCreated)
}

func (h *AccountCommandHandler) HandleGetBalance(ctx context.Context, command domain.GetAccountBalance) (decimal.Decimal, error) {
	account, err := h.accountQuery.GetAccountById(ctx, command.AccountId)

	if err != nil {
		return decimal.Zero, err
	}

	if *account.CustomerId != command.CustomerId {
		return decimal.Zero, errors.New("account doesn't belongs to user")
	}

	return h.transactionQuery.GetBalance(ctx, command.AccountId), nil
}

func NewAccountCommandHandler(
	writeRepository repositories.AccountRepository,
	accountQuery queries.AccountQuery,
	transactionQuery queries.TransactionQuery,
) AccountCommandHandler {
	return AccountCommandHandler{
		accountRepository: writeRepository,
		accountQuery:      accountQuery,
		transactionQuery:  transactionQuery,
	}
}
