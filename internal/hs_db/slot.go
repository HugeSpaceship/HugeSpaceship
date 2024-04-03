package hs_db

import (
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/model/lbp_xml/slot"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

const insertSQL = `INSERT INTO slots (
                   name, description, 
                   icon, root_level, location_x, location_y, 
                   initially_locked, sub_level, lbp1only, 
                   shareable, background, level_type, 
                   min_players, max_players, move_required, domain, uploader, first_published, last_updated, game, published)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,true) RETURNING id;`

func InsertSlot(conn *pgxpool.Conn, slot *slot.Upload, uploader uuid.UUID, game uint, domain uint) (uint64, error) {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		er2 := tx.Rollback(context.Background())
		if er2 != nil {
			return 0, er2
		}
		return 0, err
	}

	var id uint64
	err = tx.QueryRow(context.Background(), insertSQL, slot.Name, slot.Description, slot.Icon, slot.RootLevel, slot.Location.X, slot.Location.Y,
		slot.InitiallyLocked, slot.IsSubLevel, slot.IsLBP1Only, slot.Shareable,
		slot.Background, slot.LevelType, slot.MinPlayers, slot.MaxPlayers, slot.MoveRequired, domain, uploader, time.Now(), time.Now(), game,
	).Scan(&id)
	if err != nil {
		er2 := tx.Rollback(context.Background())
		if er2 != nil {
			return 0, er2
		}
		return 0, err
	}

	for _, res := range slot.Resources {
		_, err := tx.Exec(context.Background(), "INSERT INTO slot_resources VALUES($1, $2)", id, res)
		if err != nil {
			log.Debug().Err(err).Str("hash", res).Msg("failed to insert slot resource")
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}
	return id, err
}

const getSlotXML = `
SELECT 
  s.id, s.uploader, s.location_x, s.location_y, s.name, s.description, s.root_level, --FIXME: Add game to hs_db schema
  s.icon, s.initially_locked, s.sub_level, s.lbp1only, s.background, s.shareable, s.min_players, s.max_players,
  s.first_published, s.last_updated, s.game,
 
  COUNT(DISTINCT h.owner) AS heart_count,
  COUNT(DISTINCT tu.owner ) AS thumbs_up_count,
  COUNT(DISTINCT td.owner ) AS thumbs_down_count,
  COUNT(DISTINCT p.main_player) AS play_count
FROM 
  slots AS s
LEFT JOIN 
  hearts AS h ON s.id = h.slot_id 
LEFT JOIN 
  thumbs AS tu ON s.id = tu.slot_id AND NOT tu.down
LEFT JOIN
  thumbs AS td ON s.id = td.slot_id AND td.down
LEFT JOIN
  scoreboard AS p ON s.id = p.slot_id
WHERE 
  (s.id = $1) AND s.published
GROUP BY
  s.id;`

func GetSlot(conn *pgxpool.Conn, id uint64) (slot.Slot, error) {

	rows, err := conn.Query(context.Background(), getSlotXML, id)
	if err != nil {
		return slot.Slot{}, err
	}
	dbSlot, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[slot.Slot])
	if err != nil {
		return slot.Slot{}, err
	}

	dbSlot.NpHandle, err = NpHandleByUserID(conn, dbSlot.UploaderID)
	if err != nil {
		return slot.Slot{}, err
	}

	dbSlot.FirstPublishedXML = strconv.FormatInt(dbSlot.FirstPublished.UnixMilli(), 10)
	dbSlot.LastUpdatedXML = strconv.FormatInt(dbSlot.LastUpdated.UnixMilli(), 10)
	dbSlot.Type = "user"
	dbSlot.Game = 1
	dbSlot.Location = common.Location{
		X: dbSlot.LocationX,
		Y: dbSlot.LocationY,
	}
	dbSlot.PublishedIn = common.GameType(strings.ToLower(string(common.LBP2)))
	dbSlot.Icon = strings.TrimSpace(dbSlot.Icon)

	return dbSlot, nil
}

func GetTotalSlots(conn pgx.Tx) (uint64, error) {
	row := conn.QueryRow(context.Background(), "SELECT COUNT(1) FROM slots;")
	var total uint64
	return total, row.Scan(&total)
}

func GetTotalSlotsByDomain(conn pgx.Tx, domain uint) (uint64, error) {
	row := conn.QueryRow(context.Background(), "SELECT COUNT(1) FROM slots WHERE domain = $1;", domain)
	var total uint64
	return total, row.Scan(&total)
}

func GetLevelOwner(conn *pgxpool.Conn, id uint64) (uploader uuid.UUID, err error) {

	row := conn.QueryRow(context.Background(), "SELECT uploader FROM slots WHERE id = $1;", id)

	err = row.Scan(&uploader)
	return
}

func DeleteSlot(conn *pgxpool.Conn, id uint64) (err error) {

	_, err = conn.Exec(context.Background(), "UPDATE slots SET published = false WHERE id = $1;", id)
	return

}

const updateSQL = `UPDATE slots SET name = $2, description = $3, 
                   icon = $4, root_level = $5, location_x = $6, location_y = $7, 
                   initially_locked = $8, sub_level = $9, lbp1only = $10, 
                   shareable = $11, background = $12, level_type = $13, 
                   min_players = $14, max_players = $15, move_required = $16, last_updated = $17, published = true
				WHERE id = $1;
`

func UpdateSlot(conn *pgxpool.Conn, slot *slot.Upload) error {

	tx, err := conn.Begin(context.Background())
	if err != nil {
		er2 := tx.Rollback(context.Background())
		if er2 != nil {
			return er2
		}
		return err
	}
	_, err = tx.Exec(context.Background(), updateSQL, slot.ID, slot.Name, slot.Description, slot.Icon, slot.RootLevel, slot.Location.X, slot.Location.Y,
		slot.InitiallyLocked, slot.IsSubLevel, slot.IsLBP1Only, slot.Shareable,
		slot.Background, slot.LevelType, slot.MinPlayers, slot.MaxPlayers, slot.MoveRequired, time.Now(),
	)
	if err != nil {
		er2 := tx.Rollback(context.Background())
		if er2 != nil {
			return er2
		}
		return err
	}

	for _, res := range slot.Resources {
		_, err := tx.Exec(context.Background(), "INSERT INTO slot_resources VALUES($1, $2)", slot.ID, res)
		if err != nil {
			log.Debug().Err(err).Str("hash", res).Msg("failed to insert slot resource")
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return err
}
