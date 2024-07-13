package slot_filter

type LatestFilter struct {
	*baseSlotFilter
}

const latestSQL = "OFFSET ? LIMIT ?;"

func (f LatestFilter) GetQueryBase() string {
	query := slotSQL + " " + latestSQL
	return query
}

func NewLatestFilter() LatestFilter {
	filter := LatestFilter{}
	filter.baseSlotFilter = newBaseSlotFilter(filter.GetQueryBase, nil)
	return filter
}
