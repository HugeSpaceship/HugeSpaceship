package slots

import (
	"HugeSpaceship/internal/api/game_api/utils"
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/hs_db/query_builder"
	"HugeSpaceship/internal/hs_db/query_builder/query_types/slot_filter"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/pkg/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSlotsByHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		userID, err := hs_db.UserIDByName(dbCtx, ctx.Query("u"))
		pageData, err := lbp_xml.GetPageinationData(ctx)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		slots, err := query_builder.RunQuery(dbCtx, slot_filter.NewSlotsByFilter(userID), pageData)
		if err != nil {
			ctx.Error(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.Render(200, utils.LBPXML{Data: slots})
	}
}

func GetSlotsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		query_builder.RenderQuery(dbCtx, ctx, slot_filter.NewLatestFilter())
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

		slot, err := hs_db.GetSlot(dbCtx, uint64(levelID))
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(404)
			return
		}

		ctx.XML(200, slot)

	}
}

func LuckyDipHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		pageData, err := lbp_xml.GetPageinationData(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
		}

		seed, err := strconv.ParseUint(ctx.Query("seed"), 10, 64)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(400)
			return
		}

		slots, err := query_builder.RunQuery(dbCtx, slot_filter.NewLuckyDipFilter(seed), pageData)
		if err != nil {
			ctx.Error(err)
		}
		ctx.XML(200, &slots)
	}
}

func HighestRatedLevelsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		pageData, err := lbp_xml.GetPageinationData(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
		}

		slots, err := query_builder.RunQuery(dbCtx, slot_filter.NewHighestRatedFilter(), pageData)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		ctx.Render(200, utils.LBPXML{Data: slots})
	}
}
