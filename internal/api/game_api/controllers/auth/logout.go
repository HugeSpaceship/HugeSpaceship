package auth

import (
	db2 "HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LogoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, _ := ctx.Get("session")

		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		err := db2.RemoveSession(dbCtx, session.(auth.Session).Token)
		if err != nil {
			er2 := ctx.Error(err)
			if er2 != nil {
				log.Error().Err(er2).Msg("Failed to push error to the errors stack")
			}
		}
	}
}
