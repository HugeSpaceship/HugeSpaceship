package middlewares

import (
	"context"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/auth"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/utils"
	"net/http"
)

func TicketAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("MM_AUTH")
		if err != nil {
			utils.HttpLog(w, http.StatusUnauthorized, "Authentication token required")
			return
		}

		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}
		session, exists := auth.GetSession(conn, token.Value)

		if !exists {
			utils.HttpLog(w, http.StatusForbidden, "Invalid user")
			return
		}

		var domain uint
		switch session.Platform {
		case common.PS3, common.PS4, common.RPCS3:
			domain = 0
		case common.PSVita:
			domain = 1
		case common.PSP:
			domain = 2
		}

		ctx := context.WithValue(r.Context(), "domain", domain)
		ctx = context.WithValue(ctx, "session", session)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
