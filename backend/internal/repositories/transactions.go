package repositories

import (
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"sync"
)

type TransactionWriteRepository interface {
	SaveDeposit(created domain.DepositedTransaction) error
	SaveWithdraw(created domain.WithdrawnTransaction) error
}

type TransactionReadRepository interface {
	GetBalance(accountId uuid.UUID) decimal.Decimal
}

type dummyTransactionWriteRepository struct {
	mux  sync.RWMutex
	list []any
}

func (r *dummyTransactionWriteRepository) SaveDeposit(created domain.DepositedTransaction) error {
	return r.save(created)
}

func (r *dummyTransactionWriteRepository) SaveWithdraw(created domain.WithdrawnTransaction) error {
	return r.save(created)
}

func (r *dummyTransactionWriteRepository) save(event any) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.list = append(r.list, event)

	return nil
}

func (r *dummyTransactionWriteRepository) GetBalance(accountId uuid.UUID) decimal.Decimal {
	r.mux.RLock()
	defer r.mux.RUnlock()

	amount := decimal.Zero
	for _, event := range r.list {
		switch t := event.(type) {
		case domain.DepositedTransaction:
			id := t.Transaction.AccountId
			if *id == accountId {
				amount = amount.Add(t.Transaction.Amount)
			}

		case domain.WithdrawnTransaction:
			id := t.Transaction.AccountId
			if *id == accountId {
				amount = amount.Sub(t.Transaction.Amount)
			}
		}
	}

	return amount
}

func NewDummyTransactionWriteRepository() TransactionWriteRepository {
	r := &dummyTransactionWriteRepository{
		list: []any{},
	}

	return r
}

func ToDummyTransactionReadRepository(repository TransactionWriteRepository) TransactionReadRepository {
	switch r := repository.(type) {
	case *dummyTransactionWriteRepository:
		return r
	default:
		return nil
	}
}
