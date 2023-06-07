package controllers

import (
	"HugeSpaceship/pkg/common/model/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EulaHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "INSERT LICENSE TEXT HERE")
	}
}

func AnnounceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := c.Get("session")
		c.String(http.StatusOK, "Welcome to hell %s", session.(auth.Session).Username)
	}
}
