package auth

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LogoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, _ := ctx.Get("session")
		err := db.GetConnection().RemoveSession(session.(auth.Session).Token)
		if err != nil {
			er2 := ctx.Error(err)
			if er2 != nil {
				log.Error().Err(er2).Msg("Failed to push error to the errors stack")
			}
		}
	}
}
