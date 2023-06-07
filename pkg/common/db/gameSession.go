package db

import (
	"HugeSpaceship/pkg/common/model"
	"net/netip"
)

const CreateSQL = `INSERT INTO sessions (userId, ip, token, game, platform) VALUES ($1,$2,$3,$4,$5)`

func (c *Context) NewSession(username string, gameType model.GameType, ip netip.Addr, platform model.Platform, token string) error {
	userId, err := c.GetUserID(username)
	if err != nil {
		return nil

	}

	_, err = c.pool.Exec(c.ctx, CreateSQL, userId, ip, token, gameType, platform)
	return err
}

func (c *Context) GetUserID(username string) (int, error) {
	row := c.pool.QueryRow(c.ctx, "SELECT id FROM users WHERE username = $1 LIMIT 1;", username) // there can be only one

	var id int
	err := row.Scan(&id)

	return id, err
}

//func GetUser(username string) (model.User, error) {
//
//}
//
//// IsBanned checks to see if the user specified by userId is banned.
//// banned is true if the user is banned, false otherwise.
//// err is returned if something went wrong talking to the DB, for instance if the user doesn't exist.
//func IsBanned(userId int) (banned bool, err error) {
//
//}
