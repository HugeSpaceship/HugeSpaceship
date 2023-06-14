package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *Context) GetConfigSection(section string) (pgtype.Hstore, error) {
	row := c.pool.QueryRow(c.ctx, "SELECT values FROM config WHERE section = $1", section)
	store := pgtype.Hstore{}
	err := row.Scan(&store)
	return store, err
}

var storeSQL = `
IF EXISTS(SELECT FROM config WHERE section = $1) THEN
	UPDATE config SET values = $2 WHERE section = $1;
ELSE
	INSERT INTO config VALUES($1, $2);
END
`

func (c *Context) StoreConfig(section string, store pgtype.Hstore) error {
	_, err := c.pool.Exec(c.ctx, storeSQL, section, store)
	if err != nil {
		return err
	}
	return nil
}
