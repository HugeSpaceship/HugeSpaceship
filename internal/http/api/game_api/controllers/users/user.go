package users

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/auth"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/utils"
	"log/slog"
	"net/http"
	"strconv"
)

func UserGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		session := utils.GetContextValue[auth.Session](r.Context(), "session")

		username := r.PathValue("username")
		user, err := db.GetUserByName(conn, username, session.Game)
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

		user.Game = strconv.Itoa(session.Game.ToInt())

		err = utils.XMLMarshal(w, user)
		if err != nil {
			slog.Error("Failed to send response", slog.Any("error", err))
		}
	}
}
