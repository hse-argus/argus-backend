package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func EnableCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "https://argus.appweb.space")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.Writer.WriteHeader(http.StatusOK)
		return
	}

	c.Next()
}
