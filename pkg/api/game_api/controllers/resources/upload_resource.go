package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func UploadResources() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bytes, _ := io.ReadAll(ctx.Request.Body)
		err := os.WriteFile(ctx.Param("hash"), bytes, 0644)
		if err != nil {
			log.Error().Err(err).Msg("Failed to write file")
		}
	}
}
