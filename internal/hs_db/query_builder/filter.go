package query_builder

import (
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"github.com/jackc/pgx/v5"
)

type SearchFilter interface {
	RunQuery(tx pgx.Tx, domain int, skip, take uint) (slot.PaginatedSlotList[slot.SearchSlot], error)
	GetQueryBase() string
}
