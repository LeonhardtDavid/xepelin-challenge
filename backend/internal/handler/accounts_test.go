package handler

import (
	"context"
	"errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func Test_HandleCreate_NoError(t *testing.T) {
	ctx := context.Background()
	accountRepository := testAccountRepository{}
	accountQuery := testAccountQuery{}
	transactionQuery := testTransactionQuery{}
	handler := NewAccountCommandHandler(&accountRepository, &accountQuery, &transactionQuery)

	accountId := uuid.New()
	customerId := uuid.New()
	expectedAccount := domain.Account{
		Id:            &accountId,
		Name:          "account name",
		AccountNumber: "123456789",
		CustomerId:    &customerId,
	}

	err := handler.HandleCreate(ctx, domain.CreateAccount{
		Id:      uuid.New(),
		Account: expectedAccount,
		Time:    time.Now(),
	})

	assert.Nil(t, err)
	assert.Equal(t, 1, int(accountRepository.acc.Load()))
	assert.Equal(t, 0, int(transactionQuery.acc.Load()))
	assert.Equal(t, expectedAccount, accountRepository.event.Account)
}

func Test_HandleCreate_RepositoryError(t *testing.T) {
	ctx := context.Background()
	accountRepository := testErrorAccountRepository{}
	accountQuery := testAccountQuery{}
	transactionQuery := testTransactionQuery{}
	handler := NewAccountCommandHandler(&accountRepository, &accountQuery, &transactionQuery)

	accountId := uuid.New()
	customerId := uuid.New()
	expectedAccount := domain.Account{
		Id:            &accountId,
		Name:          "account name",
		AccountNumber: "123456789",
		CustomerId:    &customerId,
	}

	err := handler.HandleCreate(ctx, domain.CreateAccount{
		Id:      uuid.New(),
		Account: expectedAccount,
		Time:    time.Now(),
	})

	assert.NotNil(t, err)
	assert.Equal(t, accountRepository.err, err)
	assert.Equal(t, 0, int(transactionQuery.acc.Load()))
}

func Test_HandleGetBalance_NoError(t *testing.T) {
	ctx := context.Background()

	accountId := uuid.New()
	customerId := uuid.New()
	accountRepository := testAccountRepository{}
	accountQuery := testAccountQuery{
		account: domain.Account{
			Id:            &accountId,
			Name:          "Some Name",
			AccountNumber: "1234567890",
			CustomerId:    &customerId,
		},
	}
	transactionQuery := testTransactionQuery{}

	handler := NewAccountCommandHandler(&accountRepository, &accountQuery, &transactionQuery)

	balance, err := handler.HandleGetBalance(ctx, domain.GetAccountBalance{
		Id:         uuid.New(),
		AccountId:  accountId,
		CustomerId: customerId,
	})

	assert.Nil(t, err)
	assert.Equal(t, 0, int(accountRepository.acc.Load()))
	assert.Equal(t, 1, int(transactionQuery.acc.Load()))
	assert.Equal(t, transactionQuery.balance, balance)
}

func Test_HandleGetBalance_AccountDoesntBelongToCustomer(t *testing.T) {
	ctx := context.Background()

	accountId := uuid.New()
	customerId := uuid.New()
	anotherCustomerId := uuid.New()
	accountRepository := testAccountRepository{}
	accountQuery := testAccountQuery{
		account: domain.Account{
			Id:            &accountId,
			Name:          "Some Name",
			AccountNumber: "1234567890",
			CustomerId:    &anotherCustomerId,
		},
	}
	transactionQuery := testTransactionQuery{}

	handler := NewAccountCommandHandler(&accountRepository, &accountQuery, &transactionQuery)

	_, err := handler.HandleGetBalance(ctx, domain.GetAccountBalance{
		Id:         uuid.New(),
		AccountId:  accountId,
		CustomerId: customerId,
	})

	assert.NotNil(t, err)
	assert.Equal(t, 0, int(accountRepository.acc.Load()))
	assert.Equal(t, 0, int(transactionQuery.acc.Load()))
	assert.Equal(t, errors.New("account doesn't belongs to user"), err)
}

func Test_HandleGetBalance_ErrorGettingAccount(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("some error getting account")

	accountId := uuid.New()
	customerId := uuid.New()
	accountRepository := testAccountRepository{}
	accountQuery := testErrorAccountQuery{
		err: expectedErr,
	}
	transactionQuery := testTransactionQuery{}

	handler := NewAccountCommandHandler(&accountRepository, &accountQuery, &transactionQuery)

	_, err := handler.HandleGetBalance(ctx, domain.GetAccountBalance{
		Id:         uuid.New(),
		AccountId:  accountId,
		CustomerId: customerId,
	})

	assert.NotNil(t, err)
	assert.Equal(t, 0, int(accountRepository.acc.Load()))
	assert.Equal(t, 0, int(transactionQuery.acc.Load()))
	assert.Equal(t, expectedErr, err)
}

type testAccountRepository struct {
	acc   atomic.Int64
	event domain.AccountCreated
}

func (r *testAccountRepository) Save(_ context.Context, created domain.AccountCreated) error {
	r.acc.Add(1)
	r.event = created
	return nil
}

type testErrorAccountRepository struct {
	err error
}

type testAccountQuery struct {
	account domain.Account
}

func (r *testAccountQuery) GetAccountById(_ context.Context, _ uuid.UUID) (*domain.Account, error) {
	return &r.account, nil
}

type testErrorAccountQuery struct {
	err error
}

func (r *testErrorAccountQuery) GetAccountById(_ context.Context, _ uuid.UUID) (*domain.Account, error) {
	return nil, r.err
}

func (r *testErrorAccountRepository) Save(_ context.Context, _ domain.AccountCreated) error {
	r.err = errors.New("some error accessing repository")
	return r.err
}

type testTransactionQuery struct {
	acc       atomic.Int64
	accountId uuid.UUID
	balance   decimal.Decimal
}

func (q *testTransactionQuery) GetBalance(_ context.Context, accountId uuid.UUID) (decimal.Decimal, error) {
	q.acc.Add(1)
	q.accountId = accountId
	q.balance = decimal.NewFromFloat(rand.Float64())
	return q.balance, nil
}
