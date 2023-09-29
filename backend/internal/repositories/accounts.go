package repositories

import (
	"context"
	"encoding/json"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	Save(ctx context.Context, created domain.AccountCreated) error
}

type postgresAccountRepository struct {
	dbpool *pgxpool.Pool
}

func (r *postgresAccountRepository) Save(ctx context.Context, created domain.AccountCreated) error {
	tx, err := r.dbpool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	account := created.Account
	_, err = tx.Exec(
		ctx,
		`INSERT INTO accounts (account_id, name, account_number, customer_id) VALUES ($1, $2, $3, $4)`,
		account.Id, account.Name, account.AccountNumber, account.CustomerId,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	bytes, err := json.Marshal(&account)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	_, err = tx.Exec(
		ctx,
		`INSERT INTO account_logs (account_log_id, account, time) VALUES ($1, $2, $3)`,
		created.Id, string(bytes), created.Time,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func NewPostgresAccountRepository(dbpool *pgxpool.Pool) AccountRepository {
	return &postgresAccountRepository{
		dbpool: dbpool,
	}
}
