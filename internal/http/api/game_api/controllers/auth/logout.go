package auth

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	db2 "github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/auth"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"log/slog"
	"net/http"
)

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := utils.GetContextValue[auth.Session](r.Context(), "session")

		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		err = db2.RemoveSession(conn, session.Token)
		if err != nil {
			utils.HttpLog(w, 500, "failed to log out")
			slog.Error("failed to log out", "error", err, "userID", session.UserID)
		}
	}
}
