package photos

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/utils"
	"log/slog"
	"net/http"
	"strconv"
)

func GetPhotosBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		user, err := db.GetUserByName(conn, r.URL.Query().Get("user"), common.LBP2)
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
		photos, err := db.GetPhotos(conn, user.ID, pageSize, pageStart, domain)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "failed to get photos")
			return
		}

		err = utils.XMLMarshal(w, photos)
		if err != nil {
			slog.Error("Failed to marshal XML", slog.Any("error", err))
			return
		}
	}
}
