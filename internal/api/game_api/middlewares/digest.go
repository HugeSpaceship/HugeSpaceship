package middlewares

import (
	digest "HugeSpaceship/internal/api/game_api/utils"
	"HugeSpaceship/internal/config"
	"HugeSpaceship/pkg/utils"
	"bytes"
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
)

// Header names for ps3/vita digests

const DigestHeaderA = "X-Digest-A"
const DigestHeaderB = "X-Digest-B"

// DeferredWriter is a fairly simple gin ResponseWriter that doesn't write all its headers immediately.
// This allows a digest to be written even after the digest middlewares has executed.
// Probably not the best way of doing it, but it works so whatever.
type DeferredWriter struct {
	http.ResponseWriter
	authCookie      string
	clientDigest    string
	path            string
	alternateDigest bool
	cfg             *config.Config
}

// Write expands on the normal ResponseWriter functionality by adding digest calculation to it
func (w DeferredWriter) Write(data []byte) (int, error) {
	digestKey := w.cfg.API.DigestKey
	if w.alternateDigest {
		digestKey = w.cfg.API.AlternateDigestKey
	}
	w.Header().Add(DigestHeaderA, digest.CalculateDigest(w.path, w.authCookie, digestKey, data, false))
	return w.ResponseWriter.Write(data)
}

// NewDeferredWriter is a utility function that creates a DeferredWriter.
func NewDeferredWriter(writer http.ResponseWriter, path, clientDigest, authCookie string, alternateDigest bool, cfg *config.Config) DeferredWriter {
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
	d := digest.CalculateDigest(path, cookie, cfg.API.DigestKey, body, excludeBody)

	alternateDigest := false

	if d != digestHeader && digestHeader != "" {
		d = digest.CalculateDigest(path, cookie, cfg.API.AlternateDigestKey, body, excludeBody)
		alternateDigest = true
		if d != digestHeader {
			if cfg.API.EnforceDigest {
				return "", alternateDigest, errors.New("invalid digest")
			}
		}
	}

	return d, alternateDigest, nil
}

// DigestMiddleware calculates the digests that LBP 1 & 3 expect, this involves hashing several values from the request.
// /upload is handled differently because of the file sizes involved, this is because normally the body is hashed.
func DigestMiddleware(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			digestHeader := DigestHeaderA
			excludeBody := false
			// Check if we're on the upload endpoint because that expects things to be done differently
			if strings.Contains(r.URL.Path, "/upload/") {
				digestHeader = DigestHeaderB
				excludeBody = true
			}

			cookie, _ := r.Cookie("MM_AUTH") // if the cookie doesn't exist then we continue anyway

			var body []byte
			if !excludeBody {
				body, _ = io.ReadAll(r.Body) // Read the body of the request, unless it's a resource on /upload.
			}

			// Get the digest of the request
			d, altDigest, err := GetRequestDigest(cfg, r.URL.Path, r.Header.Get(digestHeader), cookie.Value, body, excludeBody)
			if cfg.API.EnforceDigest && err != nil { // Digest failed to authenticated
				log.Info().Str("client", r.RemoteAddr).Msg("Failed to authenticate digest, aborting request")
				utils.HttpLog(w, http.StatusForbidden, "Failed to authenticate digest.")
			} else if err != nil {
				log.Debug().Str("client", r.RemoteAddr).Msg("Digest failed to authenticate but we have digest enforcement switched off")
			}

			// Set up the writer that's used for digest verified requests
			w.Header().Set(DigestHeaderB, d)
			deferredWriter := NewDeferredWriter(w, r.URL.Path, d, cookie.Value, altDigest, cfg)

			if !excludeBody { // Re-add the request body in case anything downstream needs to use it
				r.Body = io.NopCloser(bytes.NewReader(body))
			}

			if excludeBody { // If we're on /upload then just calculate the digest here and now based on the request.
				r.Header.Set(DigestHeaderA, digest.CalculateDigest(r.URL.Path, cookie.Value, d, nil, true))
			}

			next.ServeHTTP(deferredWriter, r) // Go to the next handler in the middlewares chain
		}
		return http.HandlerFunc(fn)
	}

}
