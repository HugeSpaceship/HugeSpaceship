package resources

import (
	db2 "HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/pkg/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadResources() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, _ := ctx.Get("session")

		dbCtx := db.GetContext()

		exists, err := db2.ResourceExists(dbCtx, ctx.Param("hash"))
		if err != nil {
			ctx.Error(err)
			ctx.String(http.StatusInternalServerError, "Failed to check if resource exists")
			return
		}

		//oi im back but im not also uh i cant find the ps3! :(
		if exists {
			ctx.String(http.StatusConflict, "Resource already exists")
			return
		}

		err = db2.UploadResource(dbCtx, ctx.Request.Body, ctx.Request.ContentLength, ctx.Param("hash"), session.(auth.Session).UserID)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(500)
			return
		}

		ctx.Status(200)
	}
}
