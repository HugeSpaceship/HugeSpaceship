package db

import "github.com/google/uuid"

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
