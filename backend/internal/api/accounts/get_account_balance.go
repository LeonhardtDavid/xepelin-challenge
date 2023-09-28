package accounts

import (
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func GetBalance(handler handler.AccountCommandHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := parseAndValidate(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}

		balance, err := handler.HandleGetBalance(
			domain.GetAccountBalance{
				Id:        uuid.New(),
				AccountId: *id,
			},
		)
		if err != nil {
			ctx.Error(&errors.BadRequestError{
				Message: fmt.Sprintf("Account Id %s not found", id),
			})
			return
		}

		ctx.JSON(http.StatusOK, domain.Balance{
			AccountId: *id,
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
