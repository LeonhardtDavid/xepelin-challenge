package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	CustomerIdHeader = "X-Customer-Id" // TODO Set header using an API Gateway that handles the authentication
)

func GetCustomerId(ctx *gin.Context) uuid.UUID {
	customerId, _ := ctx.Get(CustomerIdHeader) // It's always present at this point, it's handled by middleware.RetrieveCustomer
	return customerId.(uuid.UUID)
}
