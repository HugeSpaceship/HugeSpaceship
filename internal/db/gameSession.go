package db

import (
	"context"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/auth"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const CreateSQL = `INSERT INTO sessions (userId, ip, token, game, platform, expiry) VALUES ($1,$2,$3,$4,$5,$6);`

func GetGameFromSession(session auth.Session) uint {
	switch session.Game {
	case common.LBP1, common.LBPPSP:
		return 0
	case common.LBP2, common.LBPV:
		return 1
	case common.LBP3:
		return 2
	default:
		return 0
	}
}

func NewSession(conn *pgxpool.Conn, username string, gameType common.GameType, ip netip.Addr, platform common.Platform, token string, expiry time.Time) (auth.Session, error) {
	userID, err := GetUserID(conn, username)
	if err != nil {
		return auth.Session{}, nil

	}
	n, err := PurgeSessions(conn, userID, gameType, platform)
	if err != nil {
		return auth.Session{}, err
	}
	log.Debug().Int("clearedSessions", n).Str("user", username).Msg("Purged old sessions for user")

	_, err = conn.Exec(context.Background(), CreateSQL, userID, ip, token, gameType, platform, expiry)
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

func GetUserID(conn *pgxpool.Conn, username string) (uuid.NullUUID, error) {
	row := conn.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1 LIMIT 1;", username) // there can be only one

	var id uuid.NullUUID
	err := row.Scan(&id)

	return id, err
}

func GetSession(conn *pgxpool.Conn, token string) (auth.Session, error) {
	row := conn.QueryRow(context.Background(), "SELECT sessions.*, users.username FROM sessions INNER JOIN users ON users.id = sessions.userid WHERE token = $1;", token)
	session := auth.Session{}
	err := row.Scan(&session.UserID, &session.IP, &session.Token, &session.Game, &session.Platform, &session.ExpiryDate, &session.Username)
	return session, err
}

func PurgeSessions(conn *pgxpool.Conn, userID uuid.NullUUID, game common.GameType, platform common.Platform) (int, error) {
	rows, err := conn.Exec(context.Background(), "DELETE FROM sessions WHERE userid = $1 AND game = $2 AND platform = $3", userID, game, platform)
	return int(rows.RowsAffected()), err
}

func RemoveSession(conn *pgxpool.Conn, token string) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM sessions WHERE token = $1;", token)
	return err
}
