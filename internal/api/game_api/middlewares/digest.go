package middlewares

import (
	"HugeSpaceship/internal/api/game_api/utils"
	"HugeSpaceship/internal/config"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
)

// Header names for ps3/vita digests

const DigestHeaderA = "X-Digest-A"
const DigestHeaderB = "X-Digest-B"

// DeferredWriter is a fairly simple gin ResponseWriter that doesn't write all its headers immediately.
// This allows a digest to be written even after the digest middlewares has executed.
// Probably not the best way of doing it, but it works so whatever.
type DeferredWriter struct {
	gin.ResponseWriter
	authCookie      string
	clientDigest    string
	path            string
	alternateDigest bool
	cfg             *config.Config
}

//func (w DeferredWriter) WriteHeaderNow() {
//	// Nop basically, we will write it manually later... For lbp reasons
//}

// Write expands on the normal ResponseWriter functionality by adding digest calculation to it
func (w DeferredWriter) Write(data []byte) (int, error) {
	digestKey := w.cfg.API.DigestKey
	if w.alternateDigest {
		digestKey = w.cfg.API.AlternateDigestKey
	}
	w.Header().Add(DigestHeaderA, utils.CalculateDigest(w.path, w.authCookie, digestKey, data, false))
	return w.ResponseWriter.Write(data)
}

// NewDeferredWriter is a utility function that creates a DeferredWriter.
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

// GetRequestDigest takes various parameters from the request and produces a digest.
func GetRequestDigest(cfg *config.Config, path, digestHeader, cookie string, body []byte, excludeBody bool) (string, bool, error) {
	digest := utils.CalculateDigest(path, cookie, cfg.API.DigestKey, body, excludeBody)

	alternateDigest := false

	if digest != digestHeader && digestHeader != "" {
		digest = utils.CalculateDigest(path, cookie, cfg.API.AlternateDigestKey, body, excludeBody)
		alternateDigest = true
		if digest != digestHeader {
			if cfg.API.EnforceDigest {
				return "", alternateDigest, errors.New("invalid digest")
			}
		}
	}

	return digest, alternateDigest, nil
}

// DigestMiddleware calculates the digests that LBP 1 & 3 expect, this involves hashing several values from the request.
// /upload is handled differently because of the file sizes involved, this is because normally the body is hashed.
func DigestMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		digestHeader := DigestHeaderA
		excludeBody := false
		// Check if we're on the upload endpoint because that expects things to be done differently
		if strings.Contains(ctx.Request.URL.Path, "/upload/") {
			digestHeader = DigestHeaderB
			excludeBody = true
		}

		cookie, _ := ctx.Cookie("MM_AUTH") // if the cookie doesn't exist then we continue anyway

		var body []byte
		if !excludeBody {
			body, _ = io.ReadAll(ctx.Request.Body) // Read the body of the request, unless it's a resource on /upload.
		}

		// Get the digest of the request
		digest, altDigest, err := GetRequestDigest(cfg, ctx.Request.URL.Path, ctx.GetHeader(digestHeader), cookie, body, excludeBody)
		if cfg.API.EnforceDigest && err != nil { // Digest failed to authenticated
			log.Info().Str("client", ctx.ClientIP()).Msg("Failed to authenticate digest, aborting request")
			ctx.AbortWithStatus(403)
		} else if err != nil {
			log.Debug().Str("client", ctx.ClientIP()).Msg("Digest failed to authenticate but we have digest enforcement switched off")
		}

		// Set up the writer that's used for digest verified requests
		ctx.Header(DigestHeaderB, digest)
		deferredWriter := NewDeferredWriter(ctx.Writer, ctx.Request.URL.Path, digest, cookie, altDigest, cfg)
		ctx.Writer = deferredWriter

		if !excludeBody { // Re-add the request body in case anything downstream needs to use it
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
		}

		if excludeBody { // If we're on /upload then just calculate the digest here and now based on the request.
			ctx.Header(DigestHeaderA, utils.CalculateDigest(ctx.Request.URL.Path, cookie, digest, nil, true))
		}

		ctx.Next() // Go to the next handler in the middlewares chain
	}
}
