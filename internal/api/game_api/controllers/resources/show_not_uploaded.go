package resources

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func ShowNotUploadedHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		r := lbp_xml.Resources{}
		err := ctx.BindXML(&r)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		// This checks to see if the resources already exist in the DB
		resourcesToUpload := make([]string, 0, len(r.Resources))
		for i := range r.Resources {
			exists, err := hs_db.ResourceExists(dbCtx, r.Resources[i])
			if err != nil {
				log.Error().Err(err).Msg("failed to check if resource exists, assuming it doesn't")
			}
			if !exists {
				resourcesToUpload = append(resourcesToUpload, r.Resources[i])
			}
		}

		ctx.XML(200, lbp_xml.Resources{
			Resources: resourcesToUpload,
		})
	}
}
