package resources

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/auth"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/utils"
	"log/slog"
	"net/http"
)

func UploadResources() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}
		session := utils.GetContextValue[auth.Session](r.Context(), "session")
		hash := r.PathValue("hash")

		exists, err := db.ResourceExists(conn, hash)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to check if resource exists")
			return
		}

		//oi im back but im not also uh i cant find the ps3! :(
		// - SyngletOxygen 2023

		if exists {
			utils.HttpLog(w, http.StatusConflict, "Resource already exists")
			return
		}

		err = db.UploadResource(conn, r.Body, r.ContentLength, hash, session.UserID)
		if err != nil {
			slog.Error("error saving resource", slog.Any("error", err))
			utils.HttpLog(w, http.StatusInternalServerError, "internal error in resource upload")
			return
		}
	}
}
