package slots

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

func StartPublishHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		s, err := utils.XMLUnmarshal[slot.Upload](r)
		if err != nil {
			slog.Error("Failed to parse xml body", slog.Any("error", err))
		}

		// This checks to see if the resources already exist in the DB

		resourcesToUpload := make([]string, 0, len(s.Resources))
		for i := range s.Resources {
			exists, err := hs_db.ResourceExists(conn, s.Resources[i])
			if err != nil {
				log.Warn().Err(err).Msg("failed to check if resource exists, assuming it doesn't")
			}
			if !exists {
				resourcesToUpload = append(resourcesToUpload, s.Resources[i])
			}
		}

		utils.XMLMarshal(w, slot.StartPublishSlotResponse{
			Resource: resourcesToUpload,
		})
	}
}

func PublishHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		slotData, err := utils.XMLUnmarshal[slot.Upload](r)
		if err != nil {
			utils.HttpLog(w, http.StatusBadRequest, "invalid request body")
			return
		}
		domain := utils.GetContextValue[uint](r.Context(), "domain")
		session := utils.GetContextValue[auth.Session](r.Context(), "session")

		id := slotData.ID // Will be 0 if the slot is blank

		if slotData.ID == 0 { // If inserting
			id, err := hs_db.InsertSlot(conn, slotData, session.UserID, hs_db.GetGameFromSession(session), domain)
			if err != nil {
				utils.HttpLog(w, http.StatusInternalServerError, "failed to upload level")
				slog.Error("Failed to upload level", slog.Any("error", err))
				return
			}
			slog.Debug("Published Level", slog.Uint64("levelID", id), slog.String("user", session.Username))

		} else { // If updating
			uploader, _ := hs_db.GetLevelOwner(conn, id)
			if uploader != session.UserID {
				utils.HttpLog(w, http.StatusForbidden, "permission denied")
				return
			}
			err := hs_db.UpdateSlot(conn, slotData)
			if err != nil {
				utils.HttpLog(w, http.StatusInternalServerError, "failed to update level")
				slog.Error("failed to update level", slog.Any("error", err))
				return
			}
		}

		s, err := hs_db.GetSlot(conn, id)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "failed to get level")
			slog.Error("failed to get level", slog.Any("error", err))
			return
		}

		utils.XMLMarshal(w, &s)
	}
}

func UnPublishHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		session := utils.GetContextValue[auth.Session](r.Context(), "session")

		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			utils.HttpLog(w, http.StatusBadRequest, "invalid ID")
			return
		}

		uploader, _ := hs_db.GetLevelOwner(conn, id)
		if uploader != session.UserID {
			utils.HttpLog(w, http.StatusForbidden, "permission denied")
			return
		}

		err = hs_db.DeleteSlot(conn, id)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "failed to delete level")
		}
		utils.HttpLog(w, http.StatusOK, "OK")
	}
}
