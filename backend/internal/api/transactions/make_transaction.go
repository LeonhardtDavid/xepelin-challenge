package transactions

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

func Make(handler handler.TransactionCommandHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		transaction, err := parseAndValidate(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}
		transactionId := uuid.New()
		transaction.Id = &transactionId

		if transaction.TransactionType == domain.Deposit {
			err = handler.HandleDeposit(
				ctx,
				domain.CreateDepositTransaction{
					Id:          uuid.New(),
					Transaction: *transaction,
					Time:        time.Now(),
				},
			)
		} else {
			err = handler.HandleWithdraw(
				ctx,
				domain.CreateWithdrawTransaction{
					Id:          uuid.New(),
					Transaction: *transaction,
					Time:        time.Now(),
				},
			)
		}
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, transaction)
	}
}

func parseAndValidate(ctx *gin.Context) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		return nil, &errors.BadRequestError{
			Message: fmt.Sprintf("json doesn't match expected format: %v", err),
		}
	}

	if err := transaction.Validate(); err != nil {
		return nil, &errors.BadRequestError{
			Message: fmt.Sprintf("json contains invalid values: %v", err),
		}
	}

	// TODO validate against balance

	return &transaction, nil
}
