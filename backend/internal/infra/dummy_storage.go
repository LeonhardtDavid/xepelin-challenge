package infra

import (
	"errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/google/uuid"
	"sync"
)

type DummyAccountStorage struct {
	mux  sync.RWMutex
	list []domain.AccountCreated
}

func (s *DummyAccountStorage) GetById(id uuid.UUID) (*domain.AccountCreated, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	for _, accountEvent := range s.list {
		if *accountEvent.Account.Id == id {
			return &accountEvent, nil
		}
	}

	return nil, errors.New("account not found")
}

func (s *DummyAccountStorage) Save(account domain.AccountCreated) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.list = append(s.list, account)

	return nil
}

type DummyTransactionStorage struct {
	mux  sync.RWMutex
	list []domain.TransactionEvent
}

func (s *DummyTransactionStorage) GetByAccountId(id uuid.UUID) []domain.TransactionEvent {
	s.mux.RLock()
	defer s.mux.RUnlock()

	events := make([]domain.TransactionEvent, 0)

	for _, accountEvent := range s.list {
		if *accountEvent.GetTransaction().AccountId == id {
			events = append(events, accountEvent)
		}
	}

	return events
}

func (s *DummyTransactionStorage) Save(item domain.TransactionEvent) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.list = append(s.list, item)

	return nil
}
