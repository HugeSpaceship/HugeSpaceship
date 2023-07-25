package controllers

import (
	"HugeSpaceship"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EulaHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, HugeSpaceship.LicenseText)
	}
}

func AnnounceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Make this configurable via the config file, or better yet integrate with the DB for a news list
		c.String(http.StatusOK, "") // If it's an empty string then the client won't see it
	}
}
