package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
	return
}
