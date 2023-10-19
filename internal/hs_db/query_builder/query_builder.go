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
	return filter.RunQuery(tx, int(domain), data.Start, data.Size)
}

func RunWebQuery(ctx context.Context, filter SearchFilter, start, pageSize uint) ([]slot.SearchSlot, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(context.Background())
	if err != nil {
		return nil, err
	}
	slots, err := filter.RunQuery(tx, -1, start, pageSize)
	return slots.Slots, err
}