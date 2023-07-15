package slots

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSlotsByHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.XML(200, &lbp_xml.Slots{Total: 0, HintStart: 0})

	}
}

func GetSlotHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()

		levelID, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(400)
			return
		}
		slot, err := db.GetSlot(dbCtx, &levelID)
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(404)
			return
		}
		ctx.XML(200, slot)
	}
}
