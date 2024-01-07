package query_builder

import (
	"HugeSpaceship/internal/api/game_api/utils"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func RunQuery(ctx context.Context, filter SearchFilter, data lbp_xml.PaginationData) (slot.PaginatedSlotList[slot.SearchSlot], error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	tx, err := conn.Begin(ctx)
	defer tx.Rollback(context.Background())
	if err != nil {
		return slot.PaginatedSlotList[slot.SearchSlot]{}, err
	}
	return filter.RunQuery(tx, int(data.Domain), data.Start, data.Size)
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

func RenderQuery(ctx context.Context, ginCtx *gin.Context, filter SearchFilter) {
	pageData, err := lbp_xml.GetPageinationData(ginCtx)
	if err != nil {
		ginCtx.Status(http.StatusBadRequest)
		return
	}

	slots, err := RunQuery(ctx, filter, pageData)
	if err != nil {
		ginCtx.Error(err)
		ginCtx.Status(http.StatusInternalServerError)
		return
	}

	ginCtx.Render(200, utils.LBPXML{Data: slots})
}
