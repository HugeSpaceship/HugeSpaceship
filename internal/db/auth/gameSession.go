package auth

import (
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/auth"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/npticket/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/netip"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var sessionCache = make(map[string]auth.Session)

func NewSession(conn *pgxpool.Conn, ticket types.Ticket, ip netip.Addr, game string, titleID string) (string, error) {
	if !db.UserExists(conn, ticket.Username) {
		err := db.CreateUser(conn, ticket.Username, ticket.UserID)
		if err != nil {
			return "", err
		}
	}

	token := uuid.New().String()
	platform := common.PS3
	gameType := common.LBP2

	if ticket.Footer.Signatory == types.RPCNSignatoryID {
		platform = common.RPCS3
	}
	switch game {
	case "lbp-vita":
		platform = common.PSVita
		gameType = common.LBPV
	case "lbp-psp":
		platform = common.PSP
		gameType = common.LBPPSP
	default:
		if titleID == "" {
			trim := strings.Split(ticket.TitleID, "-")[1]
			titleID = strings.Split(trim, "_")[0]
		}

		if g, exists := common.GameIDs[titleID]; exists && game == "" {
			gameType = g
		} else {
			return "", fmt.Errorf("invalid game %s", titleID)
		}
	}

	session, err := db.NewSession(conn, ticket.Username, gameType, ip, platform, token, time.Now().Add(5*time.Hour))
	if err != nil {
		return "", err
	}

	sessionCache[session.Token] = session

	return token, nil
}

func GetSession(conn *pgxpool.Conn, token string) (session auth.Session, exists bool) {

	if session, exists := sessionCache[token]; exists {
		if time.Now().After(session.ExpiryDate) {
			delete(sessionCache, token)
			err := db.RemoveSession(conn, token)
			if err != nil {
				log.Error().Err(err).Msg("Failed to remove expired session")
			}
		}
		log.Debug().Msg("Using cached session")
		return session, exists
	}

	session, err := db.GetSession(conn, token)

	if err != nil {
		log.Debug().Err(err).Msg("failed to get session")
		return session, false
	}
	if time.Now().After(session.ExpiryDate) { // If the session is expired
		err := db.RemoveSession(conn, token)
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove expired session")
		}
		return auth.Session{}, false // The auth middlewares should NOT continue if the session doesn't exist
	}
	return session, true
}
