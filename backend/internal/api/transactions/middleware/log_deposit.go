package middleware

import (
	"bytes"
	"fmt"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"io"
	"log/slog"
)

type simpleTransaction struct {
	Type   domain.TransactionType `json:"transaction_type,omitempty"`
	Amount decimal.Decimal        `json:"amount,omitempty"`
}

func LogDepositsOver(amount decimal.Decimal) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyCopy := new(bytes.Buffer)
		if _, err := io.Copy(bodyCopy, ctx.Request.Body); err == nil {
			bodyBytes := bodyCopy.Bytes()
			ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			var transaction simpleTransaction
			if err := ctx.ShouldBindJSON(&transaction); err == nil &&
				transaction.Type == domain.Deposit &&
				transaction.Amount.GreaterThanOrEqual(amount) {
				slog.InfoContext(ctx, fmt.Sprintf("A new deposit over %q has been requested", amount))
			}
			ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		ctx.Next()
		// TODO maybe I should log if the transaction was successful
	}
}
