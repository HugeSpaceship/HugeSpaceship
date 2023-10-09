package middlewares

import (
	"HugeSpaceship/internal/api/game_api/utils"
	"HugeSpaceship/internal/config"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
)

type DeferredWriter struct {
	gin.ResponseWriter
	authCookie      string
	clientDigest    string
	path            string
	alternateDigest bool
	cfg             *config.Config
}

func (w DeferredWriter) WriteHeaderNow() {
	// Nop basically, we will write it manually later... For lbp reasons
}

func (w DeferredWriter) Write(data []byte) (int, error) {
	digestKey := w.cfg.API.DigestKey
	if w.alternateDigest {
		digestKey = w.cfg.API.AlternateDigestKey
	}
	w.Header().Add("X-Digest-A", utils.CalculateDigest(w.path, w.authCookie, digestKey, data, false))
	return w.ResponseWriter.Write(data)
}

func NewDeferredWriter(writer gin.ResponseWriter, path, clientDigest, authCookie string, alternateDigest bool, cfg *config.Config) DeferredWriter {
	return DeferredWriter{
		ResponseWriter:  writer,
		authCookie:      authCookie,
		clientDigest:    clientDigest,
		path:            path,
		alternateDigest: alternateDigest,
		cfg:             cfg,
	}
}

func DigestMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// If we're on PSP, set the required headers and bail
		if ctx.GetHeader("X-exe-v") != "" {
			ctx.Header("X-exe-v", ctx.GetHeader("X-exe-v"))
			ctx.Header("X-data-v", ctx.GetHeader("X-data-v"))
			return
		}

		digestHeader := "X-Digest-A"
		excludeBody := false
		if strings.Contains(ctx.Request.URL.Path, "/upload/") {
			digestHeader = "X-Digest-B"
			excludeBody = true
		}

		cookie, _ := ctx.Cookie("MM_AUTH") // if the cookie doesn't exist then we continue anyway

		var body []byte

		if !excludeBody {
			body, _ = io.ReadAll(ctx.Request.Body) // if the client has sent a broken body, the only one that will suffer is them
		}

		digest := utils.CalculateDigest(ctx.Request.URL.Path, cookie, cfg.API.DigestKey, body, excludeBody)

		alternateDigest := false

		if digest != ctx.GetHeader(digestHeader) && ctx.GetHeader(digestHeader) != "" {
			digest = utils.CalculateDigest(ctx.Request.URL.Path, cookie, cfg.API.AlternateDigestKey, body, excludeBody)
			alternateDigest = true
			if digest != ctx.GetHeader(digestHeader) {
				if cfg.API.EnforceDigest {
					log.Debug().Msg("Failed to authenticate digest, aborting request")
					ctx.AbortWithStatus(http.StatusForbidden)
				} else {
					log.Warn().Msg("Invalid digest from client, however digests are not enforced")
				}
			}
		}
		ctx.Header("X-Digest-B", digest)
		deferredWriter := NewDeferredWriter(ctx.Writer, ctx.Request.URL.Path, digest, cookie, alternateDigest, cfg)
		ctx.Writer = deferredWriter

		if !excludeBody {
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
		}

		if excludeBody { // Basically add a digest on upload
			ctx.Header("X-Digest-A", utils.CalculateDigest(ctx.Request.URL.Path, cookie, digest, nil, true))
		}

		ctx.Next()

	}
}
