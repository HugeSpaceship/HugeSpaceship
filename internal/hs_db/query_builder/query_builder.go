package query_builder

import (
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	httpUtils "HugeSpaceship/pkg/utils"
	"context"
	"encoding/xml"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func RunQuery(conn *pgxpool.Conn, filter SearchFilter, data lbp_xml.PaginationData) (slot.PaginatedSlotList[slot.SearchSlot], error) {
	tx, err := conn.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return slot.PaginatedSlotList[slot.SearchSlot]{}, err
	}
	return filter.RunQuery(tx, int(data.Domain), data.Start, data.Size)
}

func RunWebQuery(conn *pgxpool.Conn, filter SearchFilter, start, pageSize uint) ([]slot.SearchSlot, error) {
	tx, err := conn.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return nil, err
	}
	slots, err := filter.RunQuery(tx, -1, start, pageSize)
	return slots.Slots, err
}

func RenderQuery(conn *pgxpool.Conn, w http.ResponseWriter, r *http.Request, filter SearchFilter) {
	pageData, err := lbp_xml.GetPaginationData(r)
	if err != nil {
		httpUtils.HttpLog(w, http.StatusBadRequest, "Invalid pagination data")
		return
	}

	slots, err := RunQuery(conn, filter, pageData)
	if err != nil {
		httpUtils.HttpLog(w, http.StatusInternalServerError, "Failed to fetch level")
		return
	}

	slotBytes, err := xml.Marshal(slots)
	if err != nil {
		httpUtils.HttpLog(w, http.StatusInternalServerError, "Failed to marshal XML")
		return
	}
	_, err = w.Write(slotBytes)
	if err != nil {
		panic(err)
	}
}
