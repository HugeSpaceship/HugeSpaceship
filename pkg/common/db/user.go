package db

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"HugeSpaceship/pkg/common/model/lbp_xml/npdata"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var userCreateSQL = `INSERT INTO users(id, username) VALUES ($1,$2)`

func (c *Context) CreateUser(username string) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	_, err = c.pool.Exec(c.ctx, userCreateSQL, id, username)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) UserExists(username string) bool {
	row := c.pool.QueryRow(c.ctx, "SELECT COUNT(id) FROM users WHERE username = $1", username)

	var rows int
	err := row.Scan(&rows)

	if err != nil {
		return false
	}

	return rows > 0
}

func UserIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	var id uuid.UUID

	err := pgxscan.Get(ctx, conn, &id, "SELECT id FROM users WHERE username = $1", name)
	return id, err
}

func UsernameByID(ctx context.Context, id uuid.UUID) (string, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	var username string

	err := pgxscan.Get(ctx, conn, &username, "SELECT username FROM users WHERE id = $1", id)
	return username, err
}

func NpHandleByUserID(ctx context.Context, id uuid.UUID) (npdata.NpHandle, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	var npHandle npdata.NpHandle

	err := pgxscan.Get(ctx, conn, &npHandle, "SELECT username, avatar_hash FROM users WHERE id = $1", id)
	return npHandle, err
}

func GetUserByName(ctx context.Context, name string) (lbp_xml.User, error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)
	var user lbp_xml.User

	err := pgxscan.Get(ctx, conn, &user, "SELECT users.*, users.entitled_slots - COUNT(s) AS free_slots, COUNT(s) AS used_slots FROM users LEFT JOIN slots AS s ON s.uploader = users.id WHERE username = $1 GROUP BY users.id LIMIT 1;", name)
	user.Type = "user"
	user.Game = "1"
	user.NpHandle.Username = user.Username
	user.NpHandle.IconHash = user.AvatarHash
	user.Lbp1UsedSlots = 0
	user.Lbp2FreeSlots = user.FreeSlots
	user.Lbp3FreeSlots = user.FreeSlots
	user.Lbp2EntitledSlots = user.EntitledSlots
	return user, err
}
