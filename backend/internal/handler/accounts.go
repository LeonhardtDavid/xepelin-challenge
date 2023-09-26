package handler

import (
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/repositories"
	"github.com/google/uuid"
	"time"
)

type AccountCommandHandler struct {
	writeRepository repositories.AccountWriteRepository
}

func (h *AccountCommandHandler) Handle(command domain.CreateAccount) error {
	accountCreated := domain.AccountCreated{
		Id:      uuid.New(),
		Account: command.Account,
		Time:    time.Now(),
	}

	return h.writeRepository.Save(accountCreated)
}

func NewAccountCommandHandler(writeRepository repositories.AccountWriteRepository) AccountCommandHandler {
	return AccountCommandHandler{
		writeRepository: writeRepository,
	}
}
