package resources

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/pkg/db"
	"github.com/gin-gonic/gin"
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
		resourcesToUpload, err := hs_db.CheckResources(dbCtx, r.Resources)
		if err != nil {
			ctx.Error(err)
			ctx.Status(500)
			return
		}

		ctx.XML(200, lbp_xml.Resources{
			Resources: resourcesToUpload,
		})
	}
}
