package repositories

import (
	"context"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/infra"
)

type TransactionRepository interface {
	SaveDeposit(ctx context.Context, created domain.DepositedTransaction) error
	SaveWithdraw(ctx context.Context, created domain.WithdrawnTransaction) error
}

type dummyTransactionRepository struct {
	storage *infra.DummyTransactionStorage
}

func (r *dummyTransactionRepository) SaveDeposit(_ context.Context, created domain.DepositedTransaction) error {
	return r.save(&created)
}

func (r *dummyTransactionRepository) SaveWithdraw(_ context.Context, created domain.WithdrawnTransaction) error {
	return r.save(&created)
}

func (r *dummyTransactionRepository) save(event domain.TransactionEvent) error {
	return r.storage.Save(event)
}

func NewDummyTransactionRepository(storage *infra.DummyTransactionStorage) TransactionRepository {
	r := &dummyTransactionRepository{
		storage: storage,
	}

	return r
}
