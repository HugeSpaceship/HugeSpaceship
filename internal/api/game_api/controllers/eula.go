package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const eulaText = `
HugeSpaceship is licensed under the Apache License, Version 2.0;
You can obtain a copy of the licence at:
http://www.apache.org/licenses/LICENSE-2.0
`

func EulaHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte(eulaText))
	}
}

func AnnounceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Make this configurable via the config file, or better yet integrate with the DB for a news list
		c.String(http.StatusOK, "I love this lemon") // If it's an empty string then the client won't see it
	}
}
