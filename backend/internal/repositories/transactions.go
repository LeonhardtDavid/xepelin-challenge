package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type TransactionRepository interface {
	SaveDeposit(ctx context.Context, created domain.DepositedTransaction) error
	SaveWithdraw(ctx context.Context, created domain.WithdrawnTransaction) error
}

type postgresTransactionRepository struct {
	dbpool *pgxpool.Pool
}

func (r *postgresTransactionRepository) SaveDeposit(ctx context.Context, created domain.DepositedTransaction) error {
	return r.save(ctx, &created)
}

func (r *postgresTransactionRepository) SaveWithdraw(ctx context.Context, created domain.WithdrawnTransaction) error {
	return r.save(ctx, &created)
}

func (r *postgresTransactionRepository) save(ctx context.Context, event domain.TransactionEvent) error {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	transaction := event.GetTransaction()
	row := tx.QueryRow(
		ctx,
		"SELECT balance FROM accounts_balance WHERE account_id = $1 FOR UPDATE",
		transaction.AccountId,
	)
	var currentBalance decimal.Decimal
	if err := row.Scan(&currentBalance); err != nil && !errors.Is(pgx.ErrNoRows, err) {
		tx.Rollback(ctx)
		return err
	} else if errors.Is(pgx.ErrNoRows, err) {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO accounts_balance (account_id, balance) VALUES ($1, $2)`,
			transaction.AccountId, transaction.Amount,
		)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	} else {
		var newBalance decimal.Decimal
		if transaction.TransactionType == domain.Deposit {
			newBalance = currentBalance.Add(transaction.Amount)
		} else {
			newBalance = currentBalance.Sub(transaction.Amount)
		}
		_, err = tx.Exec(
			ctx,
			`UPDATE accounts_balance SET balance = $1 WHERE account_id = $2`,
			newBalance, transaction.AccountId,
		)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	bytes, err := json.Marshal(&transaction)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	_, err = tx.Exec(
		ctx,
		`INSERT INTO transaction_logs (transaction_log_id, transaction, time) VALUES ($1, $2, $3)`,
		event.GetId(), string(bytes), event.GetTime(),
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func NewPostgresTransactionRepository(dbpool *pgxpool.Pool) TransactionRepository {
	r := &postgresTransactionRepository{
		dbpool: dbpool,
	}

	return r
}
