package photos

import (
	"HugeSpaceship/pkg/common/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetPhotosBy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()

		user, err := db.GetUserByName(dbCtx, ctx.Query("user"))
		if err != nil {
			ctx.Error(err)
			ctx.String(http.StatusBadRequest, "Invalid User")
		}

		pageSize, err := strconv.ParseUint(ctx.Query("pageSize"), 10, 64)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(400)
			return
		}
		pageStart, err := strconv.ParseUint(ctx.Query("pageStart"), 10, 64)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(400)
			return
		}

		domain := ctx.GetUint("domain")
		photos, err := db.GetPhotos(ctx, user.ID, pageSize, pageStart, domain)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.XML(http.StatusOK, photos)
	}
}
