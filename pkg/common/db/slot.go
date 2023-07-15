package db

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const insertSQL = `INSERT INTO slots (
                   id, name, description, 
                   icon, root_level, locationX, locationY, 
                   initially_locked, sub_level, lbp1only, 
                   shareable, background, level_type, 
                   min_players, max_players, move_required, domain, uploader)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`

func InsertSlot(ctx context.Context, slot *lbp_xml.SlotData, uploader uuid.UUID, domain int) (*uuid.UUID, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			return nil, er2
		}
		return nil, err
	}

	id := uuid.New()

	_, err = tx.Exec(ctx, insertSQL, id, slot.Name, slot.Description, slot.Icon, slot.RootLevel,
		slot.Location.X, slot.Location.Y, slot.InitiallyLocked, slot.IsSubLevel, slot.IsLBP1Only, slot.Shareable,
		slot.Background, slot.LevelType, slot.MinPlayers, slot.MaxPlayers, slot.MoveRequired, domain, uploader,
	)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			return nil, er2
		}
		return nil, err
	}

	for _, res := range slot.Resource {
		_, err := tx.Exec(ctx, "INSERT INTO slot_resources VALUES($1, $2)", id, res)
		if err != nil {
			er2 := tx.Rollback(ctx)
			if er2 != nil {
				return nil, er2
			}
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &id, err
}

func GetSlot(ctx context.Context, id *uuid.UUID) (*lbp_xml.SlotData, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)

	slot := struct {
		lbp_xml.SlotData
		UserName string
	}{}
	err := pgxscan.Get(ctx, conn, &slot, "SELECT slots.*, u.username FROM slots LEFT JOIN users u on slots.uploader = u.id WHERE slots.id = $1 LIMIT 1;")
	if err != nil {
		return nil, err
	}

	var resources []*lbp_xml.SlotResource
	err = pgxscan.Select(ctx, conn, &resources, "SELECT * FROM slot_resources WHERE slot_id = $1", id)
	if err != nil {
		return nil, err
	}

	resourceHashes := make([]string, len(resources))
	for i, resource := range resources {
		resourceHashes[i] = resource.ResourceHash
	}
	return &slot.SlotData, nil
}

func GetSlots(ctx context.Context, by uuid.UUID) (lbp_xml.Slots, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	var slots lbp_xml.Slots
	err := pgxscan.Select(ctx, conn, &slots.Slots, "SELECT * FROM slots WHERE uploader = $1", by)
	if err != nil {
		return lbp_xml.Slots{}, err
	}
	for i, slot := range slots.Slots {
		var resources []*lbp_xml.SlotResource
		err = pgxscan.Select(ctx, conn, &resources, "SELECT * FROM slot_resources WHERE slot_id = $1", slot)
		if err != nil {
			return lbp_xml.Slots{}, err
		}
		slots.Slots[i].Resource = make([]string, len(resources))
		for i2, resource := range resources {
			slots.Slots[i].Resource[i2] = resource.ResourceHash
		}
	}
	return slots, nil
}
