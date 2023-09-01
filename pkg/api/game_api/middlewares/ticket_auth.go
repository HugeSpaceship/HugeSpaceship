package middlewares

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/db/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TicketAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("MM_AUTH") // TODO: double check that this is the right cookie?
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

		ctx.Set("session", session) // TODO: Make this function more secure and more correct, idk
	}
}
