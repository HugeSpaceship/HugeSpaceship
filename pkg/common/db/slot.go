package db

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

const insertSQL = "INSERT INTO slots ()"

func InsertSlot(ctx context.Context, slot lbp_xml.Slot) {
	conn := ctx.Value("conn").(*pgxpool.Conn)

	conn.Exec(ctx, insertSQL)

}
