package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.FullPath()).
			Dur("latency", time.Now().Sub(start)).
			Int("status", c.Writer.Status()).
			Strs("errors", c.Errors.Errors()).
			Str("client-ip", c.ClientIP()).
			Msg("Request received")
	}
}
