package db

import (
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
