package db

import (
	"HugeSpaceship/pkg/common/model"
	"HugeSpaceship/pkg/common/model/auth"
	"github.com/rs/zerolog/log"
	"net/netip"
)

const CreateSQL = `INSERT INTO sessions (userId, ip, token, game, platform) VALUES ($1,$2,$3,$4,$5)`

func (c *Context) NewSession(username string, gameType model.GameType, ip netip.Addr, platform model.Platform, token string) error {
	userID, err := c.GetUserID(username)
	if err != nil {
		return nil

	}

	n, err := c.PurgeSessions(userID, gameType, platform)
	log.Debug().Int("clearedSessions", n).Str("user", username).Msg("Purged old sessions for user")

	_, err = c.pool.Exec(c.ctx, CreateSQL, userID, ip, token, gameType, platform)
	return err
}

func (c *Context) GetUserID(username string) (int, error) {
	row := c.pool.QueryRow(c.ctx, "SELECT id FROM users WHERE username = $1 LIMIT 1;", username) // there can be only one

	var id int
	err := row.Scan(&id)

	return id, err
}

func (c *Context) GetSession(token string) (auth.Session, error) {
	row := c.pool.QueryRow(c.ctx, "SELECT sessions.*, users.username FROM sessions INNER JOIN users ON users.id = sessions.userid WHERE token = $1;", token)

	session := auth.Session{}
	err := row.Scan(&session.UserID, &session.IP, &session.Token, &session.Game, &session.Platform, &session.Username)

	return session, err
}

func (c *Context) PurgeSessions(userID int, game model.GameType, platform model.Platform) (int, error) {
	rows, err := c.pool.Exec(c.ctx, "DELETE FROM sessions WHERE userid = $1 AND game = $2 AND platform = $3", userID, game, platform)
	return int(rows.RowsAffected()), err
}

func (c *Context) RemoveSession(token string) error {
	_, err := c.pool.Exec(c.ctx, "DELETE FROM sessions WHERE token = $1;", token)
	return err
}
