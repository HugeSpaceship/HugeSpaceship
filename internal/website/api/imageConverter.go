package api

import (
	"HugeSpaceship/internal/config"
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	lbp_image2 "HugeSpaceship/pkg/utils/file_utils/lbp_image"
	"HugeSpaceship/pkg/validation"
	"bytes"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"log/slog"
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

func ImageConverterHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ok, hash := validation.IsHashValid(r.PathValue("hash"))
		if !ok {
			utils.HttpLog(w, http.StatusBadRequest, "Invalid resource hash")
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
					err := utils.ServeFile(w, resourceFile)
					if err != nil {
						slog.Warn("Failed to read file from cache", slog.Any("err", err))
					}
					return
				}
			}
		}

		resource, tx, _, err := hs_db.GetResource(dbCtx, hash)
		if err != nil {
			utils.HttpLog(w, http.StatusNotFound, "Resource not found")
			return
		}
		defer hs_db.CloseResource(resource, tx)

		buf := new(bytes.Buffer)

		decompressed, err := lbp_image2.DecompressImage(resource)
		if errors.Is(err, lbp_image2.InvalidMagicNumber) {
			utils.HttpLog(w, http.StatusUnsupportedMediaType, "Not an image")
			return
		} else if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to fetch image.")
		}

		err = lbp_image2.IMGToPNG(decompressed, buf)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to convert image")
			return
		}

		imgReader := bytes.NewReader(buf.Bytes())

		w.Header().Set("Content-Type", "image/png")
		_, err = io.Copy(w, imgReader)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to serve image")
		}

		if cfg.ResourceServer.CacheResources {
			cacheImage(cfg.ResourceServer.CacheLocation, buf.Bytes(), hash)
		}
	}
}
