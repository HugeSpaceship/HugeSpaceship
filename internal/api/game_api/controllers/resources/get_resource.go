package resources

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/pkg/db"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
	"regexp"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func GetResourceHandler(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hash := ctx.Param("hash")
		if cfg.ResourceServer.CacheResources { // check for cache
			filePath := path.Join(cfg.ResourceServer.CacheLocation, nonAlphanumericRegex.ReplaceAllString(hash, ""))
			_, err := os.Stat(filePath)
			if err != nil {
				log.Debug().Str("hash", hash).Msg("Resource not in cache...")
			} else {
				log.Debug().Str("hash", hash).Msg("Serving resource from cache")
				ctx.File(filePath)
			}
		}
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		resource, tx, size, err := hs_db.GetResource(dbCtx, hash)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(404)
			return
		}
		defer func(tx pgx.Tx, dbCtx context.Context, ctx *gin.Context) {
			err := resource.Close()
			if err != nil {
				ctx.Error(err)
			}
			err = tx.Commit(dbCtx)
			if err != nil {
				ctx.Error(err)
			}
		}(tx, dbCtx, ctx)
		ctx.DataFromReader(200, size, "application/octet-stream", resource, nil)

		// Caches resources by resetting the read pointer on the db LOB
		if cfg.ResourceServer.CacheResources {
			_, err := resource.Seek(0, io.SeekStart)
			if err != nil {
				log.Error().Err(err).Msg("Failed to reset LOB read pointer")
				return
			}
			filePath := path.Join(cfg.ResourceServer.CacheLocation, nonAlphanumericRegex.ReplaceAllString(hash, ""))
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
			defer file.Close()
			if err != nil {
				log.Error().Err(err).Msg("Failed to open file path")
				return
			}
			_, err = io.Copy(file, resource)
			if err != nil {
				log.Error().Err(err).Msg("Failed to write resource")
			}
		}
	}
}
