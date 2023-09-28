package repositories

import (
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/infra"
)

type AccountRepository interface {
	Save(created domain.AccountCreated) error
}

type dummyAccountRepository struct {
	storage *infra.DummyAccountStorage
}

func (r *dummyAccountRepository) Save(created domain.AccountCreated) error {
	return r.storage.Save(created)
}

func NewDummyAccountRepository(storage *infra.DummyAccountStorage) AccountRepository {
	r := &dummyAccountRepository{
		storage: storage,
	}

	return r
}
