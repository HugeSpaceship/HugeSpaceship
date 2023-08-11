package auth

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/common/model/common"
	"HugeSpaceship/pkg/npticket/types"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var sessionCache = make(map[string]auth.Session)

func NewSession(ticket types.Ticket, ip netip.Addr, game string) (string, error) {
	c := db.GetConnection()
	log.Debug().Msg("got connection")
	if !c.UserExists(ticket.Username) {
		err := c.CreateUser(ticket.Username)
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
	log.Debug().Str("game", game).Msg("Game name")
	if game == "lbp-vita" {
		platform = common.PSVita
		gameType = common.LBPV
	}

	session, err := c.NewSession(ticket.Username, gameType, ip, platform, token, time.Now().Add(5*time.Hour))
	if err != nil {
		return "", err
	}

	sessionCache[session.Token] = session

	return token, nil
}

func GetSession(token string) (session auth.Session, exists bool) {

	c := db.GetConnection()

	if session, exists := sessionCache[token]; exists {
		if time.Now().After(session.ExpiryDate) {
			delete(sessionCache, token)
			err := c.RemoveSession(token)
			if err != nil {
				log.Error().Err(err).Msg("Failed to remove expired session")
			}
		}
		return session, exists
	}

	session, err := c.GetSession(token)

	if err != nil {
		log.Debug().Err(err).Msg("failed to get session")
		return session, false
	}
	if time.Now().After(session.ExpiryDate) { // If the session is expired
		err := c.RemoveSession(token)
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove expired session")
		}
		return auth.Session{}, false // The auth middleware should NOT continue if the session doesn't exist
	}
	return session, true
}
