package accounts

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Create(c *gin.Context) {
	// TODO improve it to check if the application needs to be restarted
	c.Status(http.StatusAccepted)
}
