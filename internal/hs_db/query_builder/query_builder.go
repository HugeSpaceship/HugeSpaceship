package query_builder

import (
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunQuery(ctx context.Context, filter SearchFilter, data lbp_xml.PaginationData, domain uint) (slot.PaginatedSlotList[slot.SearchSlot], error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(context.Background())
	if err != nil {
		return slot.PaginatedSlotList[slot.SearchSlot]{}, err
	}
	return filter.RunQuery(tx, domain, data)
}
