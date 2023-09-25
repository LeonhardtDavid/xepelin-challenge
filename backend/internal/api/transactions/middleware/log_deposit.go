package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"log"
)

type simpleTransaction struct {
	Type   string          `json:"transaction_type,omitempty"`
	Amount decimal.Decimal `json:"amount,omitempty"`
}

func LogDepositsOver(amount decimal.Decimal) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction simpleTransaction
		if err := c.BindJSON(&transaction); err == nil && transaction.Type == "DEPOSIT" && transaction.Amount.GreaterThanOrEqual(amount) {
			log.Println("A new deposit over", amount, "has been requested") // TODO use a better logger
		}
		c.Next()
		// TODO maybe I should log if the transaction was successful
	}
}
