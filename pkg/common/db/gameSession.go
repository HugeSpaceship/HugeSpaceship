package db

import (
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/common/model/common"
	"context"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const CreateSQL = `INSERT INTO sessions (userId, ip, token, game, platform, expiry) VALUES ($1,$2,$3,$4,$5,$6);`

func (c *Context) NewSession(username string, gameType common.GameType, ip netip.Addr, platform common.Platform, token string, expiry time.Time) (auth.Session, error) {
	userID, err := c.GetUserID(username)
	if err != nil {
		return auth.Session{}, nil

	}
	n, err := c.PurgeSessions(userID, gameType, platform)
	if err != nil {
		return auth.Session{}, err
	}
	log.Debug().Int("clearedSessions", n).Str("user", username).Msg("Purged old sessions for user")

	_, err = c.pool.Exec(c.ctx, CreateSQL, userID, ip, token, gameType, platform, expiry)
	return auth.Session{
		UserID:     userID.UUID,
		Username:   username,
		Game:       gameType,
		IP:         ip,
		Token:      token,
		ExpiryDate: expiry,
		Platform:   platform,
	}, err
}

func (c *Context) GetUserID(username string) (uuid.NullUUID, error) {
	row := c.pool.QueryRow(c.ctx, "SELECT id FROM users WHERE username = $1 LIMIT 1;", username) // there can be only one

	var id uuid.NullUUID
	err := row.Scan(&id)

	return id, err
}

func (c *Context) GetSession(token string) (auth.Session, error) {
	row := c.pool.QueryRow(context.Background(), "SELECT sessions.*, users.username FROM sessions INNER JOIN users ON users.id = sessions.userid WHERE token = $1;", token)
	session := auth.Session{}
	err := row.Scan(&session.UserID, &session.IP, &session.Token, &session.Game, &session.Platform, &session.ExpiryDate, &session.Username)
	return session, err
}

func (c *Context) PurgeSessions(userID uuid.NullUUID, game common.GameType, platform common.Platform) (int, error) {
	rows, err := c.pool.Exec(c.ctx, "DELETE FROM sessions WHERE userid = $1 AND game = $2 AND platform = $3", userID, game, platform)
	return int(rows.RowsAffected()), err
}

func (c *Context) RemoveSession(token string) error {
	_, err := c.pool.Exec(c.ctx, "DELETE FROM sessions WHERE token = $1;", token)
	return err
}
