package pages

import (
	"HugeSpaceship/internal/hs_db/query_builder"
	"HugeSpaceship/internal/hs_db/query_builder/query_types/slot_filter"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/pkg/db"
	"github.com/gin-gonic/gin"
)

func HomePage(info common.Info) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		slots, err := query_builder.RunWebQuery(dbCtx, slot_filter.NewLuckyDipFilter(0), 1, 20)
		if err != nil {
			c.Error(err)
			c.String(500, "I fucked it up somehow")
			return
		}
		c.HTML(200, "home.gohtml", gin.H{
			"Info":   info,
			"Levels": slots,
		})
	}
}

func EarthPage(info common.Info) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		slots, err := query_builder.RunWebQuery(dbCtx, slot_filter.NewLuckyDipFilter(0), 1, 20)
		if err != nil {
			c.Error(err)
			c.String(500, "I fucked it up somehow")
		}
		c.HTML(200, "earth.gohtml", gin.H{
			"Info":   info,
			"Levels": slots,
		})
	}
}
