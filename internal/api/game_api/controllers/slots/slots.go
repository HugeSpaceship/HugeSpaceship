package slots

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/hs_db/query_builder"
	"HugeSpaceship/internal/hs_db/query_builder/query_types/slot_filter"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/pkg/db"
	httpUtils "HugeSpaceship/pkg/utils"
	"net/http"
	"strconv"
)

const invalidPaginationDataError = "Invalid pagination data"
const xmlMarshalError = "Failed to marshal xml"
const levelFetchError = "Failed to fetch levels"

func GetSlotsByHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		userID, err := hs_db.UserIDByName(dbCtx, r.URL.Query().Get("u"))
		pageData, err := lbp_xml.GetPaginationData(r)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusBadRequest, invalidPaginationDataError)
			return
		}

		slots, err := query_builder.RunQuery(dbCtx, slot_filter.NewSlotsByFilter(userID), pageData)
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

func GetSlotsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		query_builder.RenderQuery(dbCtx, w, r, slot_filter.NewLatestFilter())
	}
}

func GetSlotHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)
		levelID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusBadRequest, "Invalid level id")
			return
		}
		if levelID == 0 {
			httpUtils.HttpLog(w, http.StatusNotFound, "Invalid level")
		}

		slot, err := hs_db.GetSlot(dbCtx, uint64(levelID))
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
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

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

		slots, err := query_builder.RunQuery(dbCtx, slot_filter.NewLuckyDipFilter(seed), pageData)
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
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)

		pageData, err := lbp_xml.GetPaginationData(r)
		if err != nil {
			httpUtils.HttpLog(w, http.StatusBadRequest, invalidPaginationDataError)
			return
		}

		slots, err := query_builder.RunQuery(dbCtx, slot_filter.NewHighestRatedFilter(), pageData)
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
