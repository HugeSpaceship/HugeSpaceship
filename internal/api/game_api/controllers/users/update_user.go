package users

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/pkg/db"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func UpdateUserHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, _ := ctx.Get("session")
		dbCtx := db.GetContext()

		userUpdate := lbp_xml.UpdateUser{}
		planetUpdate := lbp_xml.PlanetUpdate{}

		data, err := ctx.GetRawData()
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		err = xml.Unmarshal(data, &userUpdate)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		er2 := xml.Unmarshal(data, &planetUpdate)
		if er2 != nil {
			log.Debug().Err(er2).Msg("no bueno")
		}
		if planetUpdate.Planets != "" || planetUpdate.CCPlanet != "" {
			err := hs_db.UpdatePlanet(dbCtx, session.(auth.Session).UserID, planetUpdate, session.(auth.Session).Game)
			if err != nil {
				log.Error().Err(err).Msg("Failed to update user")
				ctx.Status(http.StatusBadRequest)
				return
			}
		}

		err = hs_db.UpdateUser(dbCtx, session.(auth.Session).UserID, userUpdate)
		if err != nil {
			log.Error().Err(err).Msg("Failed to update user")
			ctx.Status(http.StatusBadRequest)
			return
		}
	}
}
