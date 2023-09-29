package queries

import (
	"context"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/infra"
	"github.com/google/uuid"
)

type AccountQuery interface {
	GetAccountById(ctx context.Context, accountId uuid.UUID) (*domain.Account, error)
}

type dummyAccountQuery struct {
	storage *infra.DummyAccountStorage
}

func (r *dummyAccountQuery) GetAccountById(_ context.Context, accountId uuid.UUID) (*domain.Account, error) {
	event, err := r.storage.GetById(accountId)
	if err != nil {
		return nil, err
	}

	return &event.Account, nil
}

func NewDummyAccountQuery(storage *infra.DummyAccountStorage) AccountQuery {
	r := &dummyAccountQuery{
		storage: storage,
	}

	return r
}
