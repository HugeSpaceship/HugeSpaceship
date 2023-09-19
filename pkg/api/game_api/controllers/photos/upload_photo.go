package photos

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/common/model/lbp_xml/photos"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadPhoto() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var photo = photos.UploadPhoto{}
		err := ctx.BindXML(&photo)
		if err != nil {
			ctx.Error(err)
			ctx.String(http.StatusBadRequest, "Bad Request")
			return
		}

		domain := ctx.GetUint("domain")
		session, _ := ctx.Get("session")

		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		photoID, err := db.InsertPhoto(dbCtx, photo, session.(auth.Session).UserID, domain)
		if err != nil {
			ctx.Error(err)
			ctx.String(http.StatusInternalServerError, "Failed to upload photo")
			return
		}

		ctx.XML(200,
			struct {
				XMLName xml.Name `xml:"photo"`
				ID      uint64   `xml:"id"`
			}{ID: photoID},
		)
	}
}
