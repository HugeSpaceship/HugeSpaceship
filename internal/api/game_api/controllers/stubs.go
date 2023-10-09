package controllers

import "github.com/gin-gonic/gin"

func StubEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(200)
	}
}
