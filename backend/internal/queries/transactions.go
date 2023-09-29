package queries

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type TransactionQuery interface {
	GetBalance(ctx context.Context, accountId uuid.UUID) (decimal.Decimal, error)
}

type postgresTransactionQuery struct {
	dbpool *pgxpool.Pool
}

func (r *postgresTransactionQuery) GetBalance(ctx context.Context, accountId uuid.UUID) (decimal.Decimal, error) {
	row := r.dbpool.QueryRow(
		ctx,
		"SELECT balance FROM accounts_balance WHERE account_id = $1",
		accountId,
	)
	var currentBalance decimal.Decimal
	if err := row.Scan(&currentBalance); err != nil && !errors.Is(pgx.ErrNoRows, err) {
		return decimal.Zero, err
	} else if errors.Is(pgx.ErrNoRows, err) {
		return decimal.Zero, nil
	}

	return currentBalance, nil
}

func NewDummyTransactionQuery(dbpool *pgxpool.Pool) TransactionQuery {
	return &postgresTransactionQuery{
		dbpool: dbpool,
	}
}
