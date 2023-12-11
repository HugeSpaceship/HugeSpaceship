package settings

import (
	"github.com/gin-gonic/gin"
	"os"
)

func ChallengeConfigHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := os.ReadFile("ChallengeConfig.xml")
		if err != nil {
			ctx.Error(err)
		}
		ctx.Data(200, "text/xml", data)
	}
}
