package slot_filter

import (
	"github.com/rs/zerolog/log"
)

type HighestRatedFilter struct {
	*baseSlotFilter
}

const highestRatedSQL = "ORDER BY thumbs_up_count OFFSET $2 LIMIT $3;"

func (f HighestRatedFilter) GetQueryBase() string {
	query := slotSQL + " " + highestRatedSQL
	log.Debug().Str("query", query).Msg("SQL query for Highest Rated")
	return query
}

func NewHighestRatedFilter() HighestRatedFilter {
	filter := HighestRatedFilter{}
	filter.baseSlotFilter = newBaseSlotFilter(filter.GetQueryBase, nil)
	return filter
}
