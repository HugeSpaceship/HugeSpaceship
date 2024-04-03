package photos

import (
	db2 "HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/internal/model/lbp_xml/photos"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"encoding/xml"
	"net/http"
)

func UploadPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		photo, err := utils.XMLUnmarshal[photos.UploadPhoto](r)
		if err != nil {
			utils.HttpLog(w, http.StatusBadRequest, "invalid XML")
			return
		}

		domain := utils.GetContextValue[uint](r.Context(), "domain")
		session := utils.GetContextValue[auth.Session](r.Context(), "session")

		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		photoID, err := db2.InsertPhoto(dbCtx, photo, session.UserID, domain)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to upload photo")
			return
		}

		utils.XMLMarshal(w, &struct {
			XMLName xml.Name `xml:"photo"`
			ID      uint64   `xml:"id"`
		}{ID: photoID})
	}
}
