package middlewares

import (
	"HugeSpaceship/internal/hs_db/auth"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"context"
	"net/http"
)

func TicketAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("MM_AUTH")
			if err != nil {
				utils.HttpLog(w, http.StatusInternalServerError, "Failed to read cookie")
				return
			}

			dbCtx := db.GetContext()
			defer db.CloseContext(dbCtx)
			session, exists := auth.GetSession(dbCtx, token.Value)

			if !exists {
				utils.HttpLog(w, http.StatusForbidden, "User does not exist")
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

			ctx := context.WithValue(context.Background(), "domain", domain)
			ctx = context.WithValue(ctx, "session", session)
			r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
