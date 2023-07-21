package auth

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model"
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/common/model/common"
	"HugeSpaceship/pkg/npticket/types"
	"net/netip"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NewSession(ticket types.Ticket, ip netip.Addr, game string) (string, error) {
	c := db.GetConnection()

	if !c.UserExists(ticket.Username) {
		err := c.CreateUser(ticket.Username)
		if err != nil {
			return "", err
		}
	}

	token := uuid.New().String()
	platform := model.PS3
	gameType := common.LBP2
	if ticket.Footer.Signatory == types.RPCNSignatoryID {
		platform = model.RPCS3
	}
	log.Debug().Str("game", game).Msg("Game name")
	if game == "lbp-vita" {
		platform = model.PSVita
		gameType = common.LBPV
	}

	err := c.NewSession(ticket.Username, gameType, ip, platform, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetSession(token string) (session auth.Session, exists bool) {

	c := db.GetConnection()
	session, err := c.GetSession(token)

	if err != nil {
		log.Debug().Err(err).Msg("failed to get session")
		return session, false
	}
	return session, true
}
