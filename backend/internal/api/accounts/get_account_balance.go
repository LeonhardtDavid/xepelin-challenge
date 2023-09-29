package accounts

import (
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func GetBalance(handler handler.AccountCommandHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountId, err := parseAndValidate(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}

		customerId := api.GetCustomerId(ctx)

		balance, err := handler.HandleGetBalance(
			ctx,
			domain.GetAccountBalance{
				Id:         uuid.New(),
				AccountId:  *accountId,
				CustomerId: customerId,
			},
		)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Customer %s is trying to access account %s that doesn't own", customerId, accountId))
			ctx.Error(&errors.BadRequestError{
				Message: fmt.Sprintf("Account Id %s not found", accountId),
			})
			return
		}

		ctx.JSON(http.StatusOK, domain.Balance{
			AccountId: *accountId,
			Amount:    balance,
		})
	}
}

func parseAndValidate(ctx *gin.Context) (*uuid.UUID, error) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return nil, &errors.BadRequestError{
			Message: "Invalid account id format",
		}
	}

	return &id, nil
}
