package slots

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func StartPublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slot := lbp_xml.Slot{}
		err := ctx.BindXML(&slot)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse xml body")
		}
		ctx.XML(200, lbp_xml.Slot{SlotData: lbp_xml.SlotData{Resource: []string{slot.RootLevel}, Type: "user"}})
	}
}

func PublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		slotData := new(lbp_xml.SlotData)
		err := ctx.BindXML(slotData)
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(400)
		}
		domain := ctx.GetInt("domain")
		session, _ := ctx.Get("session")

		id, err := db.InsertSlot(dbCtx, slotData, session.(auth.Session).UserID, domain)
		if err != nil {
			ctx.Error(err)
		}
		slot, err := db.GetSlot(dbCtx, id)
		ctx.XML(200, slot)
	}
}
