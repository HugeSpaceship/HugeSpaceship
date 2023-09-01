package resources

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model/auth"
	"github.com/gin-gonic/gin"
)

func UploadResources() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, _ := ctx.Get("session")

		dbCtx := db.GetContext()

		// TODO: Check if resouce exists!! THIS IS IMPORTANT
		err := db.UploadResource(dbCtx, ctx.Request.Body, ctx.Request.ContentLength, ctx.Param("hash"), session.(auth.Session).UserID)
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(200)
			return
		}

		ctx.Status(200)
	}
}
