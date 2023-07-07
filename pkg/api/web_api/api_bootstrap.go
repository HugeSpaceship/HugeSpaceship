package web_api

import (
	v1 "HugeSpaceship/pkg/api/web_api/v1/controllers"
	"github.com/gin-gonic/gin"
)

func APIBootstrap(group *gin.RouterGroup) {
	apiV1 := group.Group("/v1")
	apiV1.GET("/status", v1.StatusHandler())
}
