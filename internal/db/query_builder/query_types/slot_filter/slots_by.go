package slot_filter

import (
	"github.com/google/uuid"
)

type SlotsByFilter struct {
	*baseSlotFilter
}

const slotsBySQL = "ORDER BY thumbs_up_count OFFSET $2 LIMIT $3;"

func (f SlotsByFilter) GetQueryBase() string {
	query := slotSQL + " " + slotsBySQL
	return query
}

func NewSlotsByFilter(by uuid.UUID) SlotsByFilter {
	filter := SlotsByFilter{}
	filter.baseSlotFilter = newBaseSlotFilter(filter.GetQueryBase, []string{"s.uploader = $4"}, by)
	return filter
}
