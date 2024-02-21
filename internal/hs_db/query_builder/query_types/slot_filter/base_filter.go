package slot_filter

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5"
	"math"
	"strconv"
	"strings"
)

type baseSlotFilter struct {
	extraParams     []any
	whereConditions []string
	queryFunc       func() string
}

//go:embed slot.sql
var slotSQL string

func (b baseSlotFilter) GetQueryBase() string {
	return slotSQL
}

func (b baseSlotFilter) RunQuery(tx pgx.Tx, domain int, skip, take uint) (slot.PaginatedSlotList[slot.SearchSlot], error) {
	slots := slot.PaginatedSlotList[slot.SearchSlot]{}

	params := []any{
		skip - 1, take,
	}
	if domain >= 0 {
		params = append([]any{domain}, params...)
	}

	var whereString string
	domainParam := false
	if domain >= 0 {
		whereString = "WHERE (domain = $1)"
		domainParam = true
	}
	if b.whereConditions != nil {
		for _, where := range b.whereConditions {
			whereString += fmt.Sprintf(" AND (%s)", where)
		}
	}

	query := fmt.Sprintf(b.queryFunc(), whereString)
	placeHolderCount := strings.Count(query, "?")
	for i := 1; i < placeHolderCount+1; i++ {
		if domainParam {
			query = strings.Replace(query, "?", "$"+strconv.Itoa(i+1), 1)
		} else {
			query = strings.Replace(query, "?", "$"+strconv.Itoa(i), 1)
		}
	}

	if b.extraParams != nil {
		params = append(params, b.extraParams)
	}

	rows, err := tx.Query(context.Background(), query, params...)
	s, err := pgx.CollectRows[slot.SearchSlot](rows, pgx.RowToStructByNameLax[slot.SearchSlot])
	slots.Slots = s

	if err != nil {
		return slot.PaginatedSlotList[slot.SearchSlot]{}, err
	}

	for i, s := range slots.Slots {
		if s.Name == "" {
			slots.Slots[i].Name = "Untitled Level"
		}
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

	if domain >= 0 {
		slots.Total, err = hs_db.GetTotalSlotsByDomain(tx, uint(domain))
	} else {
		slots.Total, err = hs_db.GetTotalSlots(tx)
	}
	slots.HintStart = uint64(math.Min(float64(skip+take), float64(slots.Total)))

	return slots, err
}
func newBaseSlotFilter(sqlFunc func() string, whereConditions []string, extraParams ...any) *baseSlotFilter {
	return &baseSlotFilter{
		extraParams:     extraParams,
		whereConditions: whereConditions,
		queryFunc:       sqlFunc,
	}
}
