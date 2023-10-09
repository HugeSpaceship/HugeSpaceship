package users

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/pkg/db"
	"github.com/gin-gonic/gin"
)

func UserGetHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		session, _ := ctx.Get("session")

		user, err := hs_db.GetUserByName(dbCtx, ctx.Param("username"), session.(auth.Session).Game)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(500)
			return
		}
		if user.Username == "" {
			ctx.AbortWithStatus(404)
			return
		}
		ctx.XML(200, &user)

	}
}
