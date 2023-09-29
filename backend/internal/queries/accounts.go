package queries

import (
	"context"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountQuery interface {
	GetAccountById(ctx context.Context, accountId uuid.UUID) (*domain.Account, error)
}

type postgresAccountQuery struct {
	dbpool *pgxpool.Pool
}

func (q *postgresAccountQuery) GetAccountById(ctx context.Context, accountId uuid.UUID) (*domain.Account, error) {
	row := q.dbpool.QueryRow(
		ctx,
		"SELECT account_id, name, account_number, customer_id FROM accounts WHERE account_id = $1",
		accountId,
	)

	account := domain.Account{}
	err := row.Scan(&account.Id, &account.Name, &account.AccountNumber, &account.CustomerId)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func NewPostgresAccountQuery(dbpool *pgxpool.Pool) AccountQuery {
	return &postgresAccountQuery{
		dbpool: dbpool,
	}
}
