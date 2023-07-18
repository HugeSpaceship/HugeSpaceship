package db

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"HugeSpaceship/pkg/common/model/lbp_xml/slot"
	slot2 "HugeSpaceship/pkg/common/model/lbp_xml/slot"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"math"
)

const insertSQL = `INSERT INTO slots (
                   name, description, 
                   icon, root_level, locationX, locationY, 
                   initially_locked, sub_level, lbp1only, 
                   shareable, background, level_type, 
                   min_players, max_players, move_required, domain, uploader, first_published, last_updated)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19) RETURNING id;`

func InsertSlot(ctx context.Context, slot *slot.Upload, uploader uuid.UUID, domain int) (int64, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			return 0, er2
		}
		return 0, err
	}

	var id int64
	err = tx.QueryRow(ctx, insertSQL, slot.Name, slot.Description, slot.Icon, slot.RootLevel,
		slot.Location.X, slot.Location.Y, slot.InitiallyLocked, slot.IsSubLevel, slot.IsLBP1Only, slot.Shareable,
		slot.Background, slot.LevelType, slot.MinPlayers, slot.MaxPlayers, slot.MoveRequired, domain, uploader,
		slot.FirstPublished, slot.LastUpdated,
	).Scan(&id)
	if err != nil {
		er2 := tx.Rollback(ctx)
		if er2 != nil {
			return 0, er2
		}
		return 0, err
	}

	for _, res := range slot.Resource {
		_, err := tx.Exec(ctx, "INSERT INTO slot_resources VALUES($1, $2)", id, res)
		if err != nil {
			log.Debug().Err(err).Str("hash", res).Msg("failed to insert slot resource")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}
	return id, err
}

func GetSlot(ctx context.Context, id int64) (slot.Slot, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)

	slotData := slot.Upload{}
	err := pgxscan.Get(ctx, conn, &slotData, "SELECT * FROM slots WHERE slots.id = $1 LIMIT 1;", id)
	if err != nil {
		return slot.Slot{}, err
	}

	slotData.Type = "user"
	slotData.Location = slot.Location{
		X: slotData.LocationX,
		Y: slotData.LocationY,
	}

	slotData.LastUpdatedXML = slotData.LastUpdated.UnixMilli()
	slotData.FirstPublishedXML = slotData.FirstPublished.UnixMilli()
	username, err := UsernameByID(ctx, slotData.Uploader)
	if err != nil {
		return slot.Slot{}, err
	}
	slotData.NpHandle = lbp_xml.NpHandle{
		Username: username,
	}
	slotData.Game = 1
	slot := slot.Slot{
		Upload:              slotData,
		HeartCount:          0,
		ThumbsUp:            400,
		ThumbsDown:          0,
		AverageRating:       0,
		PlayerCount:         0,
		MatchingPlayers:     0,
		TeamPick:            true,
		CommentsEnabled:     false,
		ReviewsEnabled:      false,
		UserLBP1PlayCount:   0,
		UserLBP2PlayCount:   0,
		PublishedIn:         "lbp2",
		ReviewCount:         0,
		CommentCount:        0,
		PhotoCount:          0,
		AuthorPhotoCount:    0,
		PlayCount:           0,
		UniquePlayCount:     0,
		CompletionCount:     0,
		LBP1PlayCount:       0,
		LBP1CompletionCount: 0,
		LBP1UniquePlayCount: 0,
		LBP2PlayCount:       0,
		LBP2CompletionCount: 0,
		LBP2UniquePlayCount: 0,
		LBP3PlayCount:       0,
		LBP3CompletionCount: 0,
		LBP3UniquePlayCount: 0,
	}

	return slot, nil
}

func GetSlots(ctx context.Context, by uuid.UUID) (slot.Slots, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	var slots slot.Slots
	err := pgxscan.Select(ctx, conn, &slots.Slots, "SELECT * FROM slots WHERE uploader = $1", by)
	if err != nil {
		return slot.Slots{}, err
	}

	for i, slot := range slots.Slots {
		slots.Slots[i].Type = "user"
		slots.Slots[i].Location = slot2.Location{
			X: slot.LocationX,
			Y: slot.LocationY,
		}

		slots.Slots[i].LastUpdatedXML = slot.LastUpdated.Unix()
		slots.Slots[i].FirstPublishedXML = slot.FirstPublished.Unix()
		slots.Slots[i].PublishedIn = "lbp2"
		slots.Slots[i].Game = 2
		username, err := UsernameByID(ctx, slot.Uploader)
		if err != nil {
			return slot2.Slots{}, err
		}
		slots.Slots[i].NpHandle = lbp_xml.NpHandle{
			Username: username,
		}
	}
	slots.Total = len(slots.Slots)
	slots.HintStart = int(math.Ceil(float64(len(slots.Slots))))
	return slots, nil
}
