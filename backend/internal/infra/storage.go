package infra

import (
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/google/uuid"
)

type AccountStorage interface {
	GetById(id uuid.UUID) (*domain.AccountCreated, error)
	Save(item domain.AccountCreated) error
}

type TransactionStorage interface {
	GetByAccountId(id uuid.UUID) []domain.TransactionEvent
	Save(item domain.TransactionEvent) error
}
