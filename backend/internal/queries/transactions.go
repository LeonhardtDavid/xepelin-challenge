package queries

import (
	"context"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/infra"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionQuery interface {
	GetBalance(ctx context.Context, accountId uuid.UUID) decimal.Decimal
}

type dummyTransactionQuery struct {
	storage *infra.DummyTransactionStorage
}

func (r *dummyTransactionQuery) GetBalance(_ context.Context, accountId uuid.UUID) decimal.Decimal {
	amount := decimal.Zero
	for _, event := range r.storage.GetByAccountId(accountId) {
		if *event.GetTransaction().AccountId == accountId {
			amount = amount.Add(event.GetTransaction().Amount)
		}
	}

	return amount
}

func NewDummyTransactionQuery(storage *infra.DummyTransactionStorage) TransactionQuery {
	r := &dummyTransactionQuery{
		storage: storage,
	}

	return r
}
