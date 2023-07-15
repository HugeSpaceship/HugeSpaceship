package auth

import (
	"HugeSpaceship/pkg/common/model"
	"github.com/google/uuid"
	"net/netip"
	"time"
)

type Session struct {
	ID         int
	UserID     uuid.UUID
	Username   string
	Game       model.GameType
	IP         netip.Addr // This is so that we can track changes in ip and force a re-auth, helps prevent token theft.
	Token      string
	ExpiryDate time.Time
	Platform   model.Platform
}
