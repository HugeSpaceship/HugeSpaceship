package middlewares

import "github.com/gin-gonic/gin"

func ServerHeaderMiddleware(c *gin.Context) {
	c.Header("Server", "HugeSpaceship")
}
