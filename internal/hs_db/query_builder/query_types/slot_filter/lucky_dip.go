package slot_filter

import (
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"context"
	"github.com/jackc/pgx/v5"
)

type LuckyDipFilter struct {
	*baseSlotFilter
	randomSeed uint64
}

const luckyDipSQL = "ORDER BY random() OFFSET ? LIMIT ?;"

func (f LuckyDipFilter) GetQueryBase() string {
	query := slotSQL + " " + luckyDipSQL
	return query
}

func NewLuckyDipFilter(seed uint64) LuckyDipFilter {
	filter := LuckyDipFilter{randomSeed: seed}
	filter.baseSlotFilter = newBaseSlotFilter(filter.GetQueryBase, nil)
	return filter
}

func (f LuckyDipFilter) RunQuery(tx pgx.Tx, domain int, skip, take uint) (slot.PaginatedSlotList[slot.SearchSlot], error) {
	_, err := tx.Exec(context.Background(), "SELECT setseed($1);", 0.01*float64(f.randomSeed))
	if err != nil {
		return slot.PaginatedSlotList[slot.SearchSlot]{}, err
	}
	return f.baseSlotFilter.RunQuery(tx, domain, skip, take)
}
