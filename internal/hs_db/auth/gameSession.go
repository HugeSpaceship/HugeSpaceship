package auth

import (
	"HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/pkg/npticket/types"
	"context"
	"errors"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var sessionCache = make(map[string]auth.Session)

func NewSession(ctx context.Context, ticket types.Ticket, ip netip.Addr, game string, titleID string) (string, error) {
	log.Debug().Msg("got connection")
	if !hs_db.UserExists(ctx, ticket.Username) {
		err := hs_db.CreateUser(ctx, ticket.Username, ticket.UserID)
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
		if g, exists := common.GameIDs[titleID]; exists && game == "" {
			gameType = g
		} else {
			return "", errors.New("invalid game")
		}
	}

	session, err := hs_db.NewSession(ctx, ticket.Username, gameType, ip, platform, token, time.Now().Add(5*time.Hour))
	if err != nil {
		return "", err
	}

	sessionCache[session.Token] = session

	return token, nil
}

func GetSession(ctx context.Context, token string) (session auth.Session, exists bool) {

	if session, exists := sessionCache[token]; exists {
		if time.Now().After(session.ExpiryDate) {
			delete(sessionCache, token)
			err := hs_db.RemoveSession(ctx, token)
			if err != nil {
				log.Error().Err(err).Msg("Failed to remove expired session")
			}
		}
		return session, exists
	}

	session, err := hs_db.GetSession(ctx, token)

	if err != nil {
		log.Debug().Err(err).Msg("failed to get session")
		return session, false
	}
	if time.Now().After(session.ExpiryDate) { // If the session is expired
		err := hs_db.RemoveSession(ctx, token)
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove expired session")
		}
		return auth.Session{}, false // The auth middleware should NOT continue if the session doesn't exist
	}
	return session, true
}
