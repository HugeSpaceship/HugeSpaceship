package controllers

import "github.com/gin-gonic/gin"

// NotificationController is a stub due to unknown schema
func NotificationController() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "This won't be seen")
	}
}
