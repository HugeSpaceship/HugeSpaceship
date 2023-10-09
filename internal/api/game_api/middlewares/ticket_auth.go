package middlewares

import (
	"HugeSpaceship/internal/hs_db/auth"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/pkg/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TicketAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("MM_AUTH")
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)
		session, exists := auth.GetSession(dbCtx, token)

		if !exists {
			ctx.AbortWithStatus(http.StatusForbidden)
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
		ctx.Set("domain", domain)
		ctx.Set("session", session)
	}
}
