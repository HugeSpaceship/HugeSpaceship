package photos

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/pkg/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetPhotosBy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()

		user, err := hs_db.GetUserByName(dbCtx, ctx.Query("user"), common.LBP2)
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
		photos, err := hs_db.GetPhotos(ctx, user.ID, pageSize, pageStart, domain)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.XML(http.StatusOK, photos)
	}
}
