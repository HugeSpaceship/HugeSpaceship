package logger

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

// LoggingMiddleware is a Gin middleware that logs to zerolog.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		if query != "" {
			path = path + "?" + query
		}

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		reqChan := make(chan struct{}, 1)

		// So we know if there's a stuck request
		//log.Debug().
		//	Str("path", path).
		//	Str("method", c.Request.Method).
		//	Msg("Starting request")

		// Check to see if the request was slow.
		go func() {
			select {
			case <-ctx.Done():
				log.Warn().Str("path", path).Msg("Slow Request, has taken 500ms so far")
				return
			case <-reqChan:
			}
		}()

		// Process request
		c.Next()
		reqChan <- struct{}{}

		log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Dur("latency", time.Now().Sub(start)).
			Int("status", c.Writer.Status()).
			Strs("errors", c.Errors.Errors()).
			Str("client-ip", c.ClientIP()).
			Msg("Request received")
	}
}
