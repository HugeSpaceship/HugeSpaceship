package users

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"log/slog"
	"net/http"
)

func UserGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		session := utils.GetContextValue[auth.Session](r.Context(), "session")

		username := r.PathValue("username")
		user, err := hs_db.GetUserByName(dbCtx, username, session.Game)
		if err != nil {
			slog.Error("Failed to get user",
				slog.String("username", username),
				slog.String("game", string(session.Game)),
				slog.Any("error", err),
			)
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to get user")
			return
		}
		if user.Username == "" {
			utils.HttpLog(w, http.StatusNotFound, "User not found.")
			return
		}

		err = utils.XMLMarshal(w, user)
		if err != nil {
			slog.Error("Failed to send response", slog.Any("error", err))
		}
	}
}
