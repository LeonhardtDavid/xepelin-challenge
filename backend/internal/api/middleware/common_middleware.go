package middleware

import (
	"errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api"
	internalErrors "github.com/LeonhardtDavid/xepelin-challenge/backend/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error string
}

func HandleErrors(ctx *gin.Context) {
	ctx.Next()

	if ctx.Errors != nil {
		for _, err := range ctx.Errors {
			var badRequestError *internalErrors.BadRequestError
			var unauthorizedError *internalErrors.UnauthorizedError

			if errors.As(err.Err, &badRequestError) {
				ctx.AbortWithStatusJSON(
					http.StatusBadRequest,
					ErrorResponse{
						Error: badRequestError.Error(),
					},
				)
				return
			} else if errors.As(err.Err, &unauthorizedError) {
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					ErrorResponse{
						Error: unauthorizedError.Error(),
					},
				)
				return
			} else {
				slog.ErrorContext(ctx, "Unexpected", "error", err.Err)
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					ErrorResponse{
						Error: "Something went wrong, sorry",
					},
				)
				return
			}
		}
	}
}

func RetrieveCustomer(ctx *gin.Context) {
	if customerId, err := uuid.Parse(ctx.GetHeader(api.CustomerIdHeader)); err != nil {
		ctx.Error(
			&internalErrors.UnauthorizedError{
				Message: "Unauthorized",
			},
		)
		ctx.Abort()
	} else {
		ctx.Set(api.CustomerIdHeader, customerId)
		ctx.Next()
	}

}
