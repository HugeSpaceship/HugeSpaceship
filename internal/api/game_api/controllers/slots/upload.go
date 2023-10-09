package slots

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"HugeSpaceship/pkg/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func StartPublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)
		s := slot.Slot{}
		err := ctx.BindXML(&s)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse xml body")
		}

		// This checks to see if the resources already exist in the DB

		resourcesToUpload := make([]string, 0, len(s.Resources))
		for i := range s.Resources {
			exists, err := hs_db.ResourceExists(dbCtx, s.Resources[i])
			if err != nil {
				log.Warn().Err(err).Msg("failed to check if resource exists, assuming it doesn't")
			}
			if !exists {
				resourcesToUpload = append(resourcesToUpload, s.Resources[i])
			}
		}

		ctx.XML(200, slot.StartPublishSlotResponse{
			Resource: resourcesToUpload,
		})
	}
}

func PublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)
		slotData := new(slot.Upload)
		err := ctx.BindXML(slotData)
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(400)
		}
		domain := ctx.GetInt("domain")
		session, _ := ctx.Get("session")

		id, err := hs_db.InsertSlot(dbCtx, slotData, session.(auth.Session).UserID, hs_db.GetGameFromSession(session.(auth.Session)), domain)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(500)

			return
		}
		log.Debug().Uint64("levelID", id).Str("user", session.(auth.Session).Username).Msg("Published Level")
		s, err := hs_db.GetSlot(dbCtx, id)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(500)
			return
		}
		ctx.XML(200, &s)
	}
}

func UnPublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		session, _ := ctx.Get("session")

		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.String(400, "Invalid ID")
			return
		}

		uploader, err := hs_db.GetLevelOwner(dbCtx, id)
		if uploader != session.(auth.Session).UserID {
			ctx.String(http.StatusForbidden, "User does not own level")
		}

		err = hs_db.DeleteSlot(dbCtx, id)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to delete level")
		}

		ctx.Status(200)
	}
}
