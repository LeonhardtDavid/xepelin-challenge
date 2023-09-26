package accounts

import (
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	CustomerIdHeader = "X-Customer-Id" // TODO Set header using an API Gateway that handles the authentication
)

func Create(handler handler.AccountCommandHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account, err := parseAndValidate(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}
		accountId := uuid.New()
		account.Id = &accountId

		err = handler.Handle(
			domain.CreateAccount{
				Id:      uuid.New(),
				Account: *account,
				Time:    time.Now(),
			},
		)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.Status(http.StatusAccepted)
		ctx.Header("Location", fmt.Sprintf("/accounts/%s", accountId)) // TODO not actually implemented
	}
}

func parseAndValidate(ctx *gin.Context) (*domain.Account, error) {
	var account domain.Account
	if err := ctx.ShouldBindJSON(&account); err != nil {
		return nil, &errors.BadRequestError{
			Message: fmt.Sprintf("json doesn't match expected format: %v", err),
		}
	}

	if err := account.Validate(); err != nil {
		return nil, &errors.BadRequestError{
			Message: fmt.Sprintf("json contains invalid values: %v", err),
		}
	}

	customerId, _ := ctx.Get(CustomerIdHeader) // It's always present, it's handled by middleware.RetrieveCustomer
	customerUUID := customerId.(uuid.UUID)
	account.CustomerId = &customerUUID

	// TODO validate against balance

	return &account, nil
}
