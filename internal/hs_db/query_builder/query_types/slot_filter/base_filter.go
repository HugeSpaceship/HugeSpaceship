package slot_filter

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"context"
	_ "embed"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"math"
	"strings"
)

type baseSlotFilter struct {
	extraParams []any
	whereConds  []string
	queryFunc   func() string
}

//go:embed slot.sql
var slotSQL string

func (b baseSlotFilter) GetQueryBase() string {
	return slotSQL
}

func (b baseSlotFilter) RunQuery(tx pgx.Tx, domain uint, page lbp_xml.PaginationData) (slot.PaginatedSlotList[slot.SearchSlot], error) {
	slots := slot.PaginatedSlotList[slot.SearchSlot]{}

	params := []any{
		domain, page.Start - 1, page.Size,
	}

	whereString := "(domain = $1)"
	if b.whereConds != nil {
		for _, where := range b.whereConds {
			whereString += fmt.Sprintf(" AND (%s)", where)
		}
	}

	query := fmt.Sprintf(b.queryFunc(), whereString)
	err := pgxscan.Select(context.Background(), tx, &slots.Slots, query, append(params, b.extraParams...)...)

	if err != nil {
		return slot.PaginatedSlotList[slot.SearchSlot]{}, err
	}

	for i, s := range slots.Slots {
		slots.Slots[i].NPHandle, err = hs_db.NpHandleByUserID(tx, s.Uploader)
		slots.Slots[i].Type = "user"
		slots.Slots[i].Location = common.Location{
			X: s.LocationX,
			Y: s.LocationY,
		}
		slots.Slots[i].AverageRating = 3
		slots.Slots[i].SearchScore = 10000
		slots.Slots[i].Icon = strings.TrimSpace(s.Icon)
	}

	slots.Total, err = hs_db.GetTotalSlotsByDomain(tx, domain)
	slots.HintStart = uint64(math.Min(float64(page.Start+page.Size), float64(slots.Total)))

	return slots, err
}
func newBaseSlotFilter(sqlFunc func() string, whereConditions []string, extraParams ...any) *baseSlotFilter {
	return &baseSlotFilter{
		extraParams: extraParams,
		whereConds:  whereConditions,
		queryFunc:   sqlFunc,
	}
}
