package slots

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/common/model/lbp_xml/slot"
	slot2 "HugeSpaceship/pkg/common/model/lbp_xml/slot"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

func StartPublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		slot := slot.Slot{}
		err := ctx.BindXML(&slot)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse xml body")
		}

		// This checks to see if the resources already exist in the DB
		c := 0
		resourcesToUpload := make([]string, 0, len(slot.Resource))
		for i, _ := range slot.Resource {
			exists, err := db.ResourceExists(dbCtx, slot.Resource[i])
			if err != nil {
				log.Warn().Err(err).Msg("failed to check if resource exists, assuming it doesn't")
			}
			if !exists {
				resourcesToUpload[c] = slot.Resource[i]
				c++
			}
		}

		ctx.XML(200, slot2.Slot{Upload: slot2.Upload{Resource: resourcesToUpload, Type: "user"}})
	}
}

func PublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		slotData := new(slot.Upload)
		err := ctx.BindXML(slotData)
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(400)
		}
		domain := ctx.GetInt("domain")
		session, _ := ctx.Get("session")

		//TODO: Check if the level already exists and only update the last updated if it does
		slotData.FirstPublished = time.Now()
		slotData.LastUpdated = time.Now()

		id, err := db.InsertSlot(dbCtx, slotData, session.(auth.Session).UserID, domain)
		if err != nil {
			ctx.Error(err)
		}
		slot, err := db.GetSlot(dbCtx, id)
		if err != nil {
			ctx.Error(err)
		}
		ctx.XML(200, &slot)
	}
}
