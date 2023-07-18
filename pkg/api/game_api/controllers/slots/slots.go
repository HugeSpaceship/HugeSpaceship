package slots

import (
	"HugeSpaceship/pkg/common/db"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetSlotsByHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		userID, err := db.UserIDByName(dbCtx, ctx.Query("u"))
		slots, err := db.GetSlots(dbCtx, userID)
		if err != nil {
			ctx.Error(err)
		}
		ctx.XML(200, &slots)

	}
}

func GetSlotHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()

		levelID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(400)
		}
		if levelID == 0 {
			ctx.AbortWithStatus(404)
		}
		slot, err := db.GetSlot(dbCtx, levelID)
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(404)
			return
		}

		ctx.XML(200, slot)
	}
}
