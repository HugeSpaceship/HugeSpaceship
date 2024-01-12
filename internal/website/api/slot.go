package api

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"HugeSpaceship/pkg/db"
	"github.com/gin-gonic/gin"
	"mime"
	"strconv"
)

func SlotAPI(info common.Info) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accept := ctx.Request.Header.Get("Accept")
		contentType, _, err := mime.ParseMediaType(accept)
		if err != nil {
			ctx.String(400, "Bad Accept header value '%s'", accept)
			ctx.Error(err)
			return
		}

		slotID := ctx.Query("s")
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		levelID, err := strconv.ParseUint(slotID, 10, 64)
		s, err := hs_db.GetSlot(dbCtx, levelID)
		if err != nil {
			ctx.String(500, "Failed to get slot")
			ctx.Error(err)
			return
		}

		if ctx.Request.Header.Get("HX-Request") == "true" || contentType == "text/html" {
			slotAPIHTML(ctx, info, s)
		} else {
			slotAPIJson(ctx, s)
		}
	}
}

func slotAPIJson(ctx *gin.Context, s slot.Slot) {
	ctx.JSON(200, &s)
}

func slotAPIHTML(ctx *gin.Context, info common.Info, s slot.Slot) {
	err := info.InstanceTheme.Template.ExecuteTemplate(ctx.Writer, "slotCard.gohtml", gin.H{
		"Slot": s,
	})
	if err != nil {
		panic(err)
	}
}
