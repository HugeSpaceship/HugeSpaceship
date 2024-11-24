package v3

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/http/middleware"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
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

func APIBootstrap(pool *pgxpool.Pool) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(middleware.DBCtxMiddleware(pool))
		r.Get("/slots/{id}/card", func(writer http.ResponseWriter, request *http.Request) {
			i, err := strconv.ParseUint(request.PathValue("id"), 10, 64)
			if err != nil {
				utils.HttpLog(writer, http.StatusBadRequest, "Invalid ID")
			}
			conn, err := db.GetRequestConnection(request)
			if err != nil {
				utils.HttpLog(writer, http.StatusInternalServerError, "Database connection error")
				return
			}
			s, err := db.GetSlot(conn, i)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					utils.HttpLog(writer, http.StatusNotFound, "Slot not found")
				} else {
					utils.HttpLog(writer, http.StatusInternalServerError, "Database error")
					slog.Error("Error getting slot from DB", "error", err)
				}
				return
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
				utils.HttpLog(writer, http.StatusInternalServerError, "Failed to marshal slot")
				return
			}
			writer.Header().Add("Content-Type", "application/json")
			writer.Header().Add("Access-Control-Allow-Origin", "http://localhost:5173")
			writer.WriteHeader(http.StatusOK)
			writer.Write(slotJson)
		})
	}
}
