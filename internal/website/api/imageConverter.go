package api

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/file_utils/lbp_image"
	"HugeSpaceship/pkg/validation"
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

func cacheImage(location string, resource []byte, hash string) {
	filePath := path.Join(location, "png", hash+".png")

	err := os.WriteFile(filePath, resource, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write resource")
	}

}

func ImageConverterHandler(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ok, hash := validation.IsHashValid(ctx.Param("hash"))
		if !ok {
			ctx.String(400, "Invalid resource hash")
			return
		}

		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		if cfg.ResourceServer.CacheResources { // check for cache
			if resourceFile, exists := getImageFromCache(cfg, dbCtx, hash); exists {
				if exists, err := hs_db.ResourceExists(dbCtx, hash); err != nil || !exists {
					log.Debug().Str("hash", hash).Msg("Resource exists in cache but is not in DB, deleting from cache")
					err := os.Remove(path.Join(cfg.ResourceServer.CacheLocation, hash))
					if err != nil {
						panic(err)
					}
				} else {
					ctx.File(resourceFile)
					return
				}
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

		decompressed, err := lbp_image.DecompressImage(resource)
		if errors.Is(err, lbp_image.InvalidMagicNumber) {
			ctx.String(http.StatusUnsupportedMediaType, "Not an image")
			return
		} else if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to fetch image.")
		}

		err = lbp_image.IMGToPNG(decompressed, buf)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(500)
			return
		}

		imgReader := bytes.NewReader(buf.Bytes())

		ctx.DataFromReader(200, int64(buf.Len()), "image/png", imgReader, nil)

		if cfg.ResourceServer.CacheResources {
			cacheImage(cfg.ResourceServer.CacheLocation, buf.Bytes(), hash)
		}
	}
}
