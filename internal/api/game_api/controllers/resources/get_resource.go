package resources

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/validation"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
)

func getResourceFromCache(cfg *config.Config, dbCtx context.Context, hash string) (string, bool) {

	filePath := path.Join(cfg.ResourceServer.CacheLocation, hash)
	_, err := os.Stat(filePath)
	if err != nil {
		log.Debug().Str("hash", hash).Msg("Resource not in cache...")
		return "", false
	}

	if exists, err := hs_db.ResourceExists(dbCtx, hash); err != nil || !exists {
		log.Debug().Str("hash", hash).Msg("Resource exists in cache but is not in DB, deleting from cache")
		err := os.Remove(path.Join(cfg.ResourceServer.CacheLocation, hash))
		if err != nil {
			panic(err)
		}
	}

	log.Debug().Str("hash", hash).Msg("Serving resource from cache")
	return filePath, true
}

func cacheResource(location string, resource io.ReadSeekCloser, hash string) {
	_, err := resource.Seek(0, io.SeekStart)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reset LOB read pointer")
		return
	}
	filePath := path.Join(location, hash)
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
		ok, hash := validation.IsHashValid(ctx.Param("hash"))
		if !ok {
			ctx.String(400, "Invalid resource hash")
			return
		}

		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		if cfg.ResourceServer.CacheResources { // check for cache

			if resourceFile, exists := getResourceFromCache(cfg, dbCtx, hash); exists {
				ctx.File(resourceFile)
				return
			}
		}

		resource, tx, size, err := hs_db.GetResource(dbCtx, hash)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(404)
			return
		}
		defer hs_db.CloseResource(resource, tx)
		ctx.DataFromReader(200, size, "application/octet-stream", resource, nil)

		// Caches resources by resetting the read pointer on the db LOB
		if cfg.ResourceServer.CacheResources {
			cacheResource(cfg.ResourceServer.CacheLocation, resource, hash)
		}
	}
}
