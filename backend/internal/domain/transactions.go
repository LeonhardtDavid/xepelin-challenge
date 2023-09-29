package domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type TransactionType string

const (
	Deposit  TransactionType = "DEPOSIT"
	Withdraw TransactionType = "WITHDRAW"
)

type Transaction struct {
	Id              *uuid.UUID      `json:"id,omitempty"`
	AccountId       *uuid.UUID      `json:"account_id,omitempty"`
	TransactionType TransactionType `json:"transaction_type"`
	Amount          decimal.Decimal `json:"amount"`
}

func (t *Transaction) Validate() error {
	if t.AccountId == nil {
		return errors.New("account id is required")
	}

	if t.TransactionType != Deposit && t.TransactionType != Withdraw {
		return errors.New(fmt.Sprintf("transaction type must be %s or %s", Deposit, Withdraw))
	}

	if t.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than 0")
	}

	return nil
}

//// COMMANDS

type CreateDepositTransaction struct {
	Id          uuid.UUID
	Transaction Transaction
	Time        time.Time
}

type CreateWithdrawTransaction struct {
	Id          uuid.UUID
	Transaction Transaction
	Time        time.Time
}

//// EVENTS

type TransactionEvent interface {
	GetId() uuid.UUID
	GetTransaction() Transaction
	GetTime() time.Time
}

type DepositedTransaction struct {
	Id          uuid.UUID
	Transaction Transaction
	Time        time.Time
}

func (t *DepositedTransaction) GetId() uuid.UUID {
	return t.Id
}

func (t *DepositedTransaction) GetTransaction() Transaction {
	return t.Transaction
}

func (t *DepositedTransaction) GetTime() time.Time {
	return t.Time
}

type WithdrawnTransaction struct {
	Id          uuid.UUID
	Transaction Transaction
	Time        time.Time
}

func (t *WithdrawnTransaction) GetId() uuid.UUID {
	return t.Id
}

func (t *WithdrawnTransaction) GetTransaction() Transaction {
	return t.Transaction
}

func (t *WithdrawnTransaction) GetTime() time.Time {
	return t.Time
}
