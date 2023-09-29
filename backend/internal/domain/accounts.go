package domain

import (
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

type Account struct {
	Id            *uuid.UUID `json:"id,omitempty"`
	Name          string     `json:"name"`
	AccountNumber string     `json:"account_number"`
	CustomerId    *uuid.UUID `json:"customer_id,omitempty"`
}

func (a *Account) Validate() error {
	if strings.TrimSpace(a.Name) == "" {
		return errors.New("account name is required")
	}

	if strings.TrimSpace(a.AccountNumber) == "" {
		return errors.New("account number is required")
	}

	return nil
}

//// COMMANDS

type CreateAccount struct {
	Id      uuid.UUID
	Account Account
	Time    time.Time
}

type Balance struct {
	AccountId uuid.UUID       `json:"account_id"`
	Amount    decimal.Decimal `json:"amount"`
}

//// EVENTS

type AccountCreated struct {
	Id      uuid.UUID
	Account Account
	Time    time.Time
}

type GetAccountBalance struct {
	Id         uuid.UUID
	AccountId  uuid.UUID
	CustomerId uuid.UUID
}
