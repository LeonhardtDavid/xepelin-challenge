package queries

import (
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/infra"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionQuery interface {
	GetBalance(accountId uuid.UUID) decimal.Decimal
}

type dummyTransactionQuery struct {
	storage *infra.DummyTransactionStorage
}

func (r *dummyTransactionQuery) GetBalance(accountId uuid.UUID) decimal.Decimal {
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
