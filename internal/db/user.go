package db

import (
	"context"
	"errors"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml/npdata"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
)

var userCreateSQL = `INSERT INTO users(id, username, psn_uid, rpcn_uid) VALUES ($1,$2,$3,$4)`

func invalidResourceError() error {
	return errors.New("invalid resource")
}

func CreateUser(conn *pgxpool.Conn, username string, uid uint64) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), userCreateSQL, id, username, strconv.FormatUint(uid, 10), strconv.FormatUint(uid, 10))
	if err != nil {
		return err
	}
	return nil
}

func UserExists(conn *pgxpool.Conn, username string) bool {
	row := conn.QueryRow(context.Background(), "SELECT COUNT(id) FROM users WHERE username = $1", username)

	var rows int
	err := row.Scan(&rows)

	if err != nil {
		return false
	}

	return rows > 0
}

func UserIDByName(conn *pgxpool.Conn, name string) (uuid.UUID, error) {

	const idSQL = "SELECT id FROM users WHERE username = $1"
	rows, err := conn.Query(context.Background(), idSQL, name)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])

	return id, err
}

const nphandleSQL = "SELECT username, avatar_hash FROM users WHERE id = $1"

func NpHandleByUserID(conn *pgxpool.Conn, id uuid.UUID) (*npdata.NpHandle, error) {

	rows, err := conn.Query(context.Background(), nphandleSQL, id)
	if err != nil {
		return nil, err
	}
	npHandle, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[npdata.NpHandle])

	return &npHandle, err
}

func NpHandleByUserIDTx(tx pgx.Tx, id uuid.UUID) (*npdata.NpHandle, error) {

	rows, err := tx.Query(context.Background(), nphandleSQL, id)
	if err != nil {
		return nil, err
	}
	npHandle, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[npdata.NpHandle])

	return &npHandle, err

}

func GetUserByName(conn *pgxpool.Conn, name string, game common.GameType) (*lbp_xml.User, error) {

	const userSQL = `SELECT 
    					users.*, 
						users.entitled_slots - COUNT(s) AS free_slots, 
						COUNT(s) AS used_slots 
					 FROM users LEFT JOIN slots AS s ON s.uploader = users.id 
					 WHERE username = $1 GROUP BY users.id LIMIT 1;
	`

	rows, err := conn.Query(context.Background(), userSQL, name)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[lbp_xml.User])
	if err != nil {
		return nil, err
	}
	user.XMLName.Local = "user"

	user.Type = "user"
	user.Game = "1"
	user.NpHandle.Username = user.Username
	user.NpHandle.IconHash = user.AvatarHash.String
	user.Lbp1UsedSlots = 0
	user.Lbp2FreeSlots = user.FreeSlots
	user.Lbp3FreeSlots = user.FreeSlots
	user.Lbp2EntitledSlots = user.EntitledSlots
	user.Lbp3EntitledSlots = user.EntitledSlots
	user.ClientsConnected.LittleBigPlanet2 = true
	user.Location = common.Location{
		X: user.LocationX,
		Y: user.LocationY,
	}

	// TODO: offload the processing of this to the DB
	switch game {
	case common.LBP2:
		user.Planets = user.LBP2Planet
	case common.LBP3:
		user.Planets = user.LBP3Planet
	case common.LBPV:
		user.Planets = user.LBPVPlanet
	}

	return &user, err
}

func UpdatePlanet(conn *pgxpool.Conn, id uuid.UUID, update *lbp_xml.PlanetUpdate, game common.GameType) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	if strings.TrimSpace(update.Planets) != "" {
		if len(update.Planets) > 40 {
			return invalidResourceError()
		}
		switch game {
		case common.LBP2:
			_, err := tx.Exec(context.Background(), "UPDATE users SET planet_lbp2 = $1, planet_lbp3 = $1 WHERE id = $2", update.Planets, id)
			if err != nil {
				return err
			}
		case common.LBP3:
			_, err := tx.Exec(context.Background(), "UPDATE users SET planet_lbp3 = $1 WHERE id = $2", update.Planets, id)
			if err != nil {
				return err
			}
		case common.LBPV:
			_, err := tx.Exec(context.Background(), "UPDATE users SET planet_lbp_vita = $1 WHERE id = $2", update.Planets, id)
			if err != nil {
				return err
			}
		default:
			return errors.New("invalid game client for planet update")
		}
	}
	if strings.TrimSpace(update.CCPlanet) != "" {
		if len(update.CCPlanet) > 40 {
			return invalidResourceError()
		}
		_, err := tx.Exec(context.Background(), "UPDATE users SET planet_cc = $1 WHERE id = $2", update.CCPlanet, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func UpdateUser(conn *pgxpool.Conn, id uuid.UUID, update *lbp_xml.UpdateUser) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if strings.TrimSpace(update.Biography) != "" {
		if len(update.Biography) > 512 {
			return errors.New("biography too long")
		}
		_, err := tx.Exec(context.Background(), "UPDATE users SET bio = $1 WHERE id = $2", update.Biography, id)
		if err != nil {
			return err
		}
	}

	if strings.TrimSpace(update.Icon) != "" {
		if len(update.Icon) > 40 {
			return invalidResourceError()
		}
		_, err := tx.Exec(context.Background(), "UPDATE users SET avatar_hash = $1 WHERE id = $2", update.Icon, id)
		if err != nil {
			return err
		}
	}

	if strings.TrimSpace(update.BooHash) != "" {
		if len(update.BooHash) > 40 {
			return invalidResourceError()
		}
		_, err := tx.Exec(context.Background(), "UPDATE users SET boo_icon = $1 WHERE id = $2", update.BooHash, id)
		if err != nil {
			return err
		}
	}
	if strings.TrimSpace(update.MehHash) != "" {
		if len(update.MehHash) > 40 {
			return invalidResourceError()
		}
		_, err := tx.Exec(context.Background(), "UPDATE users SET meh_icon = $1 WHERE id = $2", update.MehHash, id)
		if err != nil {
			return err
		}
	}
	if strings.TrimSpace(update.YayHash) != "" {
		if len(update.YayHash) > 40 {
			return invalidResourceError()
		}
		_, err := tx.Exec(context.Background(), "UPDATE users SET yay_icon = $1 WHERE id = $2", update.YayHash, id)
		if err != nil {
			return err
		}
	}

	if update.Location != nil {
		_, err := tx.Exec(context.Background(), "UPDATE users SET location_x = $1, location_y = $2 WHERE id = $3",
			update.Location.X, update.Location.Y, id)
		if err != nil {
			return err
		}
	}

	for _, slot := range update.Slots.Slots {
		if o, err := GetLevelOwner(conn, slot.Id); err != nil || o != id {
			return errors.New("level not owned by user, or it does not exist")
		}

		_, err := tx.Exec(context.Background(), "UPDATE slots SET location_x = $1, location_y = $2 WHERE id = $3",
			slot.Location.X, slot.Location.Y, slot.Id)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	return err
}
