package resources

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadResources() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, _ := ctx.Get("session")

		dbCtx := db.GetContext()

		exists, err := db.ResourceExists(dbCtx, ctx.Param("hash"))
		if err != nil {
			ctx.Error(err)
			ctx.String(http.StatusInternalServerError, "Failed to check if resource exists")
			return
		}
		if exists {
			ctx.String(http.StatusConflict, "Resource already exists")
			return
		}

		err = db.UploadResource(dbCtx, ctx.Request.Body, ctx.Request.ContentLength, ctx.Param("hash"), session.(auth.Session).UserID)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(500)
			return
		}

		ctx.Status(200)
	}
}
