package v3

import (
	"encoding/json"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/middleware"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type testSlot struct {
	Name        string
	Author      string
	Description string
	Hearts      uint64
	Plays       uint64
	Game        common.GameType
}

func V3APIBootstrap() func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(middleware.DBCtxMiddleware)
		r.Get("/slots/{id}/card", func(writer http.ResponseWriter, request *http.Request) {
			i, err := strconv.ParseUint(request.PathValue("id"), 10, 64)
			if err != nil {
				utils.HttpLog(writer, http.StatusBadRequest, "Invalid ID")
			}
			conn, err := db.GetRequestConnection(request)
			if err != nil {
				utils.HttpLog(writer, http.StatusInternalServerError, "Database connection error")
			}
			s, err := db.GetSlot(conn, i)
			if err != nil {
				utils.HttpLog(writer, http.StatusInternalServerError, "Database error")
			}

			slot := testSlot{
				Name:        s.Name,
				Author:      s.NpHandle.Username,
				Description: s.Description,
				Hearts:      s.HeartCount,
				Plays:       s.PlayCount,
				Game:        s.PublishedIn,
			}

			slotJson, err := json.Marshal(&slot)
			if err != nil {
				utils.HttpLog(writer, http.StatusInternalServerError, "Database error")
			}
			writer.Header().Add("Content-Type", "application/json")
			writer.Header().Add("Access-Control-Allow-Origin", "http://xlocalhost:5173")
			writer.WriteHeader(http.StatusOK)
			writer.Write(slotJson)
		})
	}
}
