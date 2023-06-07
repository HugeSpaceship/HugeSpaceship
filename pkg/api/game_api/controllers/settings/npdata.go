package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
)

func NpDataEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic("lul")
		}
		log.Info().Msg(string(data)) // TODO: Fuck you lbp, your protocol sucks
		ctx.String(200, "sure buddy")
	}
}
