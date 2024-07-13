package photos

import (
	"encoding/xml"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	db2 "github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/auth"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml/photos"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
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

		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		photoID, err := db2.InsertPhoto(conn, photo, session.UserID, domain)
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
