package users

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"github.com/rs/zerolog/log"
	"net/http"
)

func UpdateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := utils.GetContextValue[auth.Session](r.Context(), "session")
		dbCtx := db.GetContext()

		userUpdate, err := utils.XMLUnmarshal[lbp_xml.UpdateUser](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		planetUpdate, er2 := utils.XMLUnmarshal[lbp_xml.PlanetUpdate](r)
		if er2 != nil {
			log.Debug().Err(er2).Msg("no bueno")
		}
		if planetUpdate.Planets != "" || planetUpdate.CCPlanet != "" {
			err := hs_db.UpdatePlanet(dbCtx, session.UserID, planetUpdate, session.Game)
			if err != nil {
				utils.HttpLog(w, http.StatusBadRequest, "failed to update user")
				return
			}
		}

		err = hs_db.UpdateUser(dbCtx, session.UserID, userUpdate)
		if err != nil {
			utils.HttpLog(w, http.StatusBadRequest, "failed to update user")
			return
		}
	}
}
