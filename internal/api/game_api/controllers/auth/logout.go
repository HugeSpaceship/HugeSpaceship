package auth

import (
	db2 "HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"github.com/rs/zerolog/log"
	"net/http"
)

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := utils.GetContextValue[auth.Session](r.Context(), "session")

		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		err := db2.RemoveSession(dbCtx, session.Token)
		if err != nil {
			utils.HttpLog(w, 500, "failed to log out")
			log.Error().Err(err).Msg("Failed to push error to the errors stack")
		}
	}
}
