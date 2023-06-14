package db

var userCreateSQL = `INSERT INTO users(username) VALUES ($1)`

func (c *Context) CreateUser(username string) error {
	_, err := c.pool.Exec(c.ctx, userCreateSQL, username)
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
