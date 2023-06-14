package middlewares

import (
	"HugeSpaceship/pkg/api/game_api/utils"
	"HugeSpaceship/pkg/common/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
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

	apiConfig := config.GetLBPAPIConfig()
	digestKey := apiConfig.PrimaryDigest
	if w.alternateDigest {
		digestKey = apiConfig.AlternateDigest
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

func DigestMiddleware(excludeBody bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, _ := ctx.Cookie("MM_AUTH") // if the cookie doesn't exist then we continue anyway

		body, _ := io.ReadAll(ctx.Request.Body) // if the client has sent a broken body, the only one that will suffer is them

		digest := utils.CalculateDigest(ctx.Request.URL.Path, cookie, config.GetLBPAPIConfig().PrimaryDigest, body, excludeBody)

		alternateDigest := false

		if digest != ctx.GetHeader("X-Digest-A") {
			digest = utils.CalculateDigest(ctx.Request.URL.Path, cookie, config.GetLBPAPIConfig().AlternateDigest, body, excludeBody)
			alternateDigest = true
			if digest != ctx.GetHeader("X-Digest-A") {
				log.Debug().Msg("Failed to authenticate digest, aborting request")
				ctx.AbortWithStatus(http.StatusForbidden)
			}
		}
		ctx.Header("X-Digest-B", digest)
		deferredWriter := NewDeferredWriter(ctx.Writer, ctx.Request.URL.Path, digest, cookie, alternateDigest)
		ctx.Writer = deferredWriter

		ctx.Next()

	}
}
