package interceptor

import (
	"github.com/gin-gonic/gin"
)

func GeneralInterceptor(c *gin.Context) {
	token := c.Query("token")
	if token == "1234" {
		c.Next()
	} else {
		c.JSON(401, gin.H{"message": "unauthorized"})
		c.Abort()
	}
}