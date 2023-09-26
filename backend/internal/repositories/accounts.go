package repositories

import (
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"sync"
)

type AccountWriteRepository interface {
	Save(created domain.AccountCreated) error
}

type dummyAccountWriteRepository struct {
	mux  sync.RWMutex
	list []domain.AccountCreated
}

func (r *dummyAccountWriteRepository) Save(created domain.AccountCreated) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.list = append(r.list, created)

	return nil
}

func NewDummyAccountWriteRepository() AccountWriteRepository {
	r := &dummyAccountWriteRepository{
		list: []domain.AccountCreated{},
	}

	return r
}
