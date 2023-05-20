package controllers

import (
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
		c.String(http.StatusOK, "Welcome to the first test of awful server")
	}
}
