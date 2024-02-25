package api

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log/slog"
	"mime"
	"net/http"
	"strconv"
)

func SlotAPI(info common.Info) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		contentType, _, err := mime.ParseMediaType(accept)
		if err != nil {
			utils.HttpLogf(w, 400, "Bad Accept header value '%s'", accept)
			slog.Debug("Failed to parse media type", slog.Any("err", err))
			return
		}

		slotID := r.URL.Query().Get("s")
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		levelID, err := strconv.ParseUint(slotID, 10, 64)
		s, err := hs_db.GetSlot(dbCtx, levelID)
		if err != nil {
			utils.HttpLog(w, 500, "Failed to get slot")
			slog.Error("failed to get slot", slog.String("id", slotID), slog.Any("err", err))
			return
		}

		if r.Header.Get("HX-Request") == "true" || contentType == "text/html" {
			slotAPIHTML(w, info, s)
		} else {
			slotAPIJson(w, s)
		}
	}
}

func slotAPIJson(w http.ResponseWriter, s slot.Slot) {
	jsonBytes, err := json.Marshal(&s)
	if err != nil {
		slog.Error("failed to marshal slot", slog.Any("err", err))
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		slog.Error("failed write to ResponseWriter")
	}
}

func slotAPIHTML(w http.ResponseWriter, info common.Info, s slot.Slot) {
	err := info.InstanceTheme.Template.ExecuteTemplate(w, "slotCard.gohtml", gin.H{
		"Slot": s,
	})
	if err != nil {
		panic(err)
	}
}
