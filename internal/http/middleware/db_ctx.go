package middleware

import (
	"context"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"log/slog"
	"net/http"
)

// DBCtxMiddleware creates a database connection for a handler to use, without the handler having to (and often screwing up)
// the management of a DB connection
func DBCtxMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.Acquire()
		if err != nil {
			slog.Error("failed to get database connection", slog.Any("error", err))
			utils.HttpLog(w, 500, "failed to acquire DB connection")
			return
		}
		defer conn.Release()
		ctx := context.WithValue(r.Context(), db.ConnCtxKey, conn)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
