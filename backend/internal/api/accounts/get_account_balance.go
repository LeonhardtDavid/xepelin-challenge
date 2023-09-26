package accounts

import (
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

		balance := handler.HandleGetBalance(
			domain.GetAccountBalance{
				Id:        uuid.New(),
				AccountId: *id,
			},
		)

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
			Message: "Invalid account id",
		}
	}

	// TODO validate account belongs to user

	return &id, nil
}
