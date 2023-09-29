package handler

import (
	"context"
	"errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func Test_HandleDeposit_NoError(t *testing.T) {
	ctx := context.Background()
	repository := testTransactionRepository{}
	handler := NewTransactionCommandHandler(&repository)

	transactionId := uuid.New()
	accountId := uuid.New()
	expectedTransaction := domain.Transaction{
		Id:              &transactionId,
		AccountId:       &accountId,
		TransactionType: domain.Deposit,
		Amount:          decimal.NewFromFloat(999.99),
	}

	err := handler.HandleDeposit(ctx, domain.CreateDepositTransaction{
		Id:          uuid.New(),
		Transaction: expectedTransaction,
		Time:        time.Now(),
	})

	assert.Nil(t, err)
	assert.Equal(t, 1, int(repository.acc.Load()))
	assert.Equal(t, expectedTransaction, repository.depositEvent.Transaction)
	assert.Nil(t, repository.withdrawEvent)
}

func Test_HandleDeposit_Error(t *testing.T) {
	ctx := context.Background()
	repository := testErrorTransactionRepository{}
	handler := NewTransactionCommandHandler(&repository)

	err := handler.HandleDeposit(ctx, domain.CreateDepositTransaction{})

	assert.NotNil(t, err)
	assert.Equal(t, repository.err, err)
}

func Test_HandleWithdraw_NoError(t *testing.T) {
	ctx := context.Background()
	repository := testTransactionRepository{}
	handler := NewTransactionCommandHandler(&repository)

	transactionId := uuid.New()
	accountId := uuid.New()
	expectedTransaction := domain.Transaction{
		Id:              &transactionId,
		AccountId:       &accountId,
		TransactionType: domain.Withdraw,
		Amount:          decimal.NewFromFloat(9.99),
	}

	err := handler.HandleWithdraw(ctx, domain.CreateWithdrawTransaction{
		Id:          uuid.New(),
		Transaction: expectedTransaction,
		Time:        time.Now(),
	})

	assert.Nil(t, err)
	assert.Equal(t, 1, int(repository.acc.Load()))
	assert.Equal(t, expectedTransaction, repository.withdrawEvent.Transaction)
	assert.Nil(t, repository.depositEvent)
}

func Test_HandleWithdraw_Error(t *testing.T) {
	ctx := context.Background()
	repository := testErrorTransactionRepository{}
	handler := NewTransactionCommandHandler(&repository)

	err := handler.HandleWithdraw(ctx, domain.CreateWithdrawTransaction{})

	assert.NotNil(t, err)
	assert.Equal(t, repository.err, err)
}

type testTransactionRepository struct {
	acc           atomic.Int64
	depositEvent  *domain.DepositedTransaction
	withdrawEvent *domain.WithdrawnTransaction
}

func (r *testTransactionRepository) SaveDeposit(_ context.Context, command domain.DepositedTransaction) error {
	r.acc.Add(1)
	r.depositEvent = &command
	return nil
}

func (r *testTransactionRepository) SaveWithdraw(_ context.Context, command domain.WithdrawnTransaction) error {
	r.acc.Add(1)
	r.withdrawEvent = &command
	return nil
}

type testErrorTransactionRepository struct {
	err error
}

func (r *testErrorTransactionRepository) SaveDeposit(_ context.Context, _ domain.DepositedTransaction) error {
	r.err = errors.New("some error trying to deposit money")
	return r.err
}

func (r *testErrorTransactionRepository) SaveWithdraw(_ context.Context, _ domain.WithdrawnTransaction) error {
	r.err = errors.New("some error trying to withdraw money")
	return r.err
}
