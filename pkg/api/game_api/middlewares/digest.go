package middlewares

import (
	"HugeSpaceship/pkg/api/game_api/utils"
	"HugeSpaceship/pkg/common/config"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
}

func (w DeferredWriter) WriteHeaderNow() {
	// Nop basically, we will write it manually later... For lbp reasons
}

func (w DeferredWriter) Write(data []byte) (int, error) {
	digestKey := viper.GetString("mainline_digest")
	if w.alternateDigest {
		digestKey = viper.GetString("vita_digest")
	}
	w.Header().Add("X-Digest-A", utils.CalculateDigest(w.path, w.authCookie, digestKey, data, false))
	return w.ResponseWriter.Write(data)
}

func NewDeferredWriter(writer gin.ResponseWriter, path, clientDigest, authCookie string, alternateDigest bool) DeferredWriter {
	return DeferredWriter{
		ResponseWriter:  writer,
		authCookie:      authCookie,
		clientDigest:    clientDigest,
		path:            path,
		alternateDigest: alternateDigest,
	}
}

func DigestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		digestHeader := "X-Digest-A"
		excludeBody := false
		if strings.HasPrefix(ctx.Request.URL.Path, "/api/LBP_XML/upload") {
			digestHeader = "X-Digest-B"
			excludeBody = true
		}

		cookie, _ := ctx.Cookie("MM_AUTH") // if the cookie doesn't exist then we continue anyway

		var body []byte

		if !excludeBody {
			body, _ = io.ReadAll(ctx.Request.Body) // if the client has sent a broken body, the only one that will suffer is them
		}

		digest := utils.CalculateDigest(ctx.Request.URL.Path, cookie, config.GetLBPAPIConfig().PrimaryDigest, body, excludeBody)

		alternateDigest := false

		if digest != ctx.GetHeader(digestHeader) {
			digest = utils.CalculateDigest(ctx.Request.URL.Path, cookie, config.GetLBPAPIConfig().AlternateDigest, body, excludeBody)
			alternateDigest = true
			if digest != ctx.GetHeader(digestHeader) {
				log.Debug().Msg("Failed to authenticate digest, aborting request")
				ctx.AbortWithStatus(http.StatusForbidden)
			}
		}
		ctx.Header("X-Digest-B", digest)
		deferredWriter := NewDeferredWriter(ctx.Writer, ctx.Request.URL.Path, digest, cookie, alternateDigest)
		ctx.Writer = deferredWriter

		if !excludeBody {
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
		}

		ctx.Next()
	}
}
