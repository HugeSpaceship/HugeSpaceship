package users

import (
	"HugeSpaceship/pkg/common/db"
	"github.com/gin-gonic/gin"
)

func UserGetHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()

		user, err := db.GetUserByName(dbCtx, ctx.Param("username"))
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
