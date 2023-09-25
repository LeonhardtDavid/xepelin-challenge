package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func Make(c *gin.Context) {
	// TODO improve it to check if the application needs to be restarted
	c.JSON(http.StatusCreated, gin.H{
		"id": uuid.New(),
	})
}
