package api

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/image_utils"
	"HugeSpaceship/pkg/validation"
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"path"
)

func getImageFromCache(cfg *config.Config, dbCtx context.Context, hash string) (string, bool) {

	filePath := path.Join(cfg.ResourceServer.CacheLocation, "png", hash+".png")
	_, err := os.Stat(filePath)
	if err != nil {
		log.Debug().Str("hash", hash).Msg("Resource not in cache...")
		return "", false
	}

	if exists, err := hs_db.ResourceExists(dbCtx, hash); err != nil || !exists {
		log.Debug().Str("hash", hash).Msg("Resource exists in cache but is not in DB, deleting from cache")
		err := os.Remove(path.Join(cfg.ResourceServer.CacheLocation, "png", hash+".png"))
		if err != nil {
			panic(err)
		}
	}

	log.Debug().Str("hash", hash).Msg("Serving resource from cache")
	return filePath, true
}

func cacheImage(location string, resource io.ReadSeeker, hash string) {
	_, err := resource.Seek(0, io.SeekStart)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reset LOB read pointer")
		return
	}
	filePath := path.Join(location, hash+".png")
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

func ImageConverterHandler(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		hash := ctx.Param("hash")
		if validation.IsHashValid(hash) {
			ctx.String(400, "Invalid resource hash")
			return
		}

		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		if exists, err := hs_db.ResourceExists(dbCtx, hash); err != nil || !exists {
			log.Debug().Str("hash", hash).Msg("Resource exists in cache but is not in DB, deleting from cache")
			err := os.Remove(path.Join(cfg.ResourceServer.CacheLocation, hash))
			if err != nil {
				panic(err)
			}
		}

		if cfg.ResourceServer.CacheResources { // check for cache

			if resourceFile, exists := getImageFromCache(cfg, dbCtx, hash); exists {
				ctx.File(resourceFile)
				return
			}
		}

		resource, tx, _, err := hs_db.GetResource(dbCtx, hash)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(404)
			return
		}
		defer hs_db.CloseResource(resource, tx)

		buf := new(bytes.Buffer)

		decompressed := image_utils.DecompressImage(resource)
		if decompressed == nil {
			ctx.String(http.StatusUnsupportedMediaType, "Not an image")
			return
		}
		err = image_utils.IMGToPNG(decompressed, buf)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(500)
			return
		}

		imgReader := bytes.NewReader(buf.Bytes())

		ctx.DataFromReader(200, int64(buf.Len()), "image/png", imgReader, nil)

		if cfg.ResourceServer.CacheResources {
			cacheImage(cfg.ResourceServer.CacheLocation, resource, hash)
		}
	}
}
