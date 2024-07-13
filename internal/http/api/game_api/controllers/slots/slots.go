package slots

import (
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/query_builder"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/query_builder/query_types/slot_filter"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml"
	httpUtils "github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

const invalidPaginationDataError = "Invalid pagination data"
const xmlMarshalError = "Failed to marshal xml"
const levelFetchError = "Failed to fetch levels"

func GetSlotsByHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		fmt.Println("Getting UserID for user " + r.URL.Query().Get("u"))

		userID, err := db.UserIDByName(conn, r.URL.Query().Get("u"))
		pageData, err := lbp_xml.GetPaginationData(r)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusBadRequest, invalidPaginationDataError)
			return
		}

		fmt.Println("Getting slots for user " + userID.String())

		slots, err := query_builder.RunQuery(conn, slot_filter.NewSlotsByFilter(userID), pageData)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusInternalServerError, levelFetchError)
			slog.Error("Failed to get levels", slog.Any("error", err))
			return
		}

		fmt.Println("marshalling levels for user " + userID.String())

		err = httpUtils.XMLMarshal(w, &slots)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusInternalServerError, xmlMarshalError)
			return
		}
	}
}

func GetSlotsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		query_builder.RenderQuery(conn, w, r, slot_filter.NewLatestFilter())
	}
}

func GetSlotHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		levelID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusBadRequest, "Invalid level id")
			return
		}
		if levelID == 0 {
			httpUtils.HttpLog(w, http.StatusNotFound, "Invalid level")
		}

		slot, err := db.GetSlot(conn, uint64(levelID))
		if err != nil {
			httpUtils.HttpLog(w, http.StatusNotFound, "Level not found")
			return
		}

		err = httpUtils.XMLMarshal(w, &slot)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusInternalServerError, xmlMarshalError)
			return
		}
	}
}

func LuckyDipHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		pageData, err := lbp_xml.GetPaginationData(r)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusBadRequest, invalidPaginationDataError)
			return
		}

		seed, err := strconv.ParseUint(r.URL.Query().Get("seed"), 10, 64)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusBadRequest, "Invalid seed value")
			return
		}

		slots, err := query_builder.RunQuery(conn, slot_filter.NewLuckyDipFilter(seed), pageData)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusInternalServerError, levelFetchError)
			return
		}

		err = httpUtils.XMLMarshal(w, &slots)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusInternalServerError, xmlMarshalError)
			return
		}
	}
}

func HighestRatedLevelsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		pageData, err := lbp_xml.GetPaginationData(r)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusBadRequest, invalidPaginationDataError)
			return
		}

		slots, err := query_builder.RunQuery(conn, slot_filter.NewHighestRatedFilter(), pageData)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusBadRequest, levelFetchError)
			return
		}

		err = httpUtils.XMLMarshal(w, &slots)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusInternalServerError, xmlMarshalError)
			return
		}
	}
}
