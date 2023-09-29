package accounts

import (
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func Create(handler handler.AccountCommandHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account, err := parseAndValidateAccount(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}
		accountId := uuid.New()
		account.Id = &accountId

		err = handler.HandleCreate(
			ctx,
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

		ctx.Header("Location", fmt.Sprintf("%s/%s", ctx.FullPath(), accountId)) // TODO not actually implemented
		ctx.JSON(http.StatusCreated, account)
	}
}

func parseAndValidateAccount(ctx *gin.Context) (*domain.Account, error) {
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

	customerId := api.GetCustomerId(ctx)
	account.CustomerId = &customerId

	return &account, nil
}
