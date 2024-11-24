package slot_filter

type HighestRatedFilter struct {
	*baseSlotFilter
}

const highestRatedSQL = "ORDER BY thumbs_up_count OFFSET $2 LIMIT $3;"

func (f HighestRatedFilter) GetQueryBase() string {
	query := slotSQL + " " + highestRatedSQL
	return query
}

func NewHighestRatedFilter() HighestRatedFilter {
	filter := HighestRatedFilter{}
	filter.baseSlotFilter = newBaseSlotFilter(filter.GetQueryBase, nil)
	return filter
}
