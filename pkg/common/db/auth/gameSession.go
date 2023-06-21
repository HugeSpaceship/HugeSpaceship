package auth

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model"
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/npticket/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/netip"
)

func NewSession(ticket types.Ticket, ip netip.Addr) string {
	c := db.GetConnection()

	if !c.UserExists(ticket.Username) {
		err := c.CreateUser(ticket.Username)
		if err != nil {
			panic(err.Error())
		}
	}

	token := uuid.New().String()

	err := c.NewSession(ticket.Username, model.LBP2, ip, model.PS3, token)
	if err != nil {
		panic(err.Error())
		return token
	}

	return token
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
