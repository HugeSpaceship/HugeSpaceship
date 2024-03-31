package photos

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"net/http"
	"strconv"
)

func GetPhotosBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbCtx := db.GetContext()

		user, err := hs_db.GetUserByName(dbCtx, r.URL.Query().Get("user"), common.LBP2)
		if err != nil {
			utils.HttpLog(w, http.StatusBadRequest, "invalid User")
		}

		pageSize, err := strconv.ParseUint(r.URL.Query().Get("pageSize"), 10, 64)
		if err != nil {
			utils.HttpLog(w, http.StatusBadRequest, "invalid pageSize")
			return
		}
		pageStart, err := strconv.ParseUint(r.URL.Query().Get("pageStart"), 10, 64)
		if err != nil {
			utils.HttpLog(w, http.StatusBadRequest, "invalid pageStart")
			return
		}

		domain := utils.GetContextValue[uint](r.Context(), "domain")
		photos, err := hs_db.GetPhotos(dbCtx, user.ID, pageSize, pageStart, domain)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "failed to get photos")
			return
		}

		utils.XMLMarshal(w, photos)
	}
}
