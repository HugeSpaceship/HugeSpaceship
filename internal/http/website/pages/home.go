package pages

import (
	"HugeSpaceship/internal/hs_db/query_builder"
	"HugeSpaceship/internal/hs_db/query_builder/query_types/slot_filter"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"net/http"
)

func HomePage(info common.Info) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		slots, err := query_builder.RunWebQuery(conn, slot_filter.NewLuckyDipFilter(0), 1, 20)
		if err != nil {
			utils.HttpLog(w, 500, "Failed to query DB")
			return
		}

		err = info.InstanceTheme.Template.ExecuteTemplate(w, "home.gohtml", utils.TmplMap{
			"Info":   info,
			"Levels": slots,
		})

		if err != nil {
			utils.HttpLog(w, 500, "Error in template...")
		}
	}
}
