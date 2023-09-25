package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Liveness(c *gin.Context) {
	// TODO improve it to check if the application needs to be restarted
	c.Status(http.StatusNoContent)
}

func Readiness(c *gin.Context) {
	// TODO improve it to check the application is ready to receive traffic
	c.Status(http.StatusNoContent)
}
