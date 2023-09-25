package accounts

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBalance(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"balance": 100,
	})
}
