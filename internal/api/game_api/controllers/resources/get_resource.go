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

func getResourceFromCache(cfg *config.Config, hash string) string {
	filePath := path.Join(cfg.ResourceServer.CacheLocation, hash)
	_, err := os.Stat(filePath)
	if err != nil {
		log.Debug().Str("hash", hash).Msg("Resource not in cache...")
	} else {
		log.Debug().Str("hash", hash).Msg("Serving resource from cache")
	}
	return filePath
}

func closeResource(resource io.ReadSeekCloser, tx pgx.Tx) {
	err := resource.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close resource")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit")
	}
}

func cacheResource(location string, resource io.ReadSeekCloser, hash string) {
	_, err := resource.Seek(0, io.SeekStart)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reset LOB read pointer")
		return
	}
	filePath := path.Join(location, nonAlphanumericRegex.ReplaceAllString(hash, ""))
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

func GetResourceHandler(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hash := ctx.Param("hash")
		if nonAlphanumericRegex.MatchString(hash) {
			ctx.String(400, "Invalid resource hash")
			return
		}

		if cfg.ResourceServer.CacheResources { // check for cache
			ctx.File(getResourceFromCache(cfg, hash))
			return
		}
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		resource, tx, size, err := hs_db.GetResource(dbCtx, hash)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(404)
			return
		}
		defer closeResource(resource, tx)
		ctx.DataFromReader(200, size, "application/octet-stream", resource, nil)

		// Caches resources by resetting the read pointer on the db LOB
		if cfg.ResourceServer.CacheResources {
			cacheResource(cfg.ResourceServer.CacheLocation, resource, hash)
		}
	}
}
