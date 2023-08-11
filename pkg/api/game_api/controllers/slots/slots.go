package slots

import (
	"HugeSpaceship/pkg/common/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSlotsByHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)
		userID, err := db.UserIDByName(dbCtx, ctx.Query("u"))
		pageCount, err := strconv.ParseInt(ctx.Query("pageSize"), 10, 64)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(400)
			return
		}
		pageStart, err := strconv.ParseInt(ctx.Query("pageStart"), 10, 64)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(400)
			return
		}

		slots, err := db.GetSlotsBy(dbCtx, userID, pageStart, pageCount)
		if err != nil {
			ctx.Error(err)
		}
		ctx.XML(200, &slots)
		db.CloseContext(dbCtx)
	}
}

func GetSlotHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)
		levelID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(400)
		}
		if levelID == 0 {
			ctx.AbortWithStatus(404)
		}
		slot, err := db.GetSlot(dbCtx, uint64(levelID))
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(404)
			return
		}

		ctx.XML(200, slot)

	}
}
