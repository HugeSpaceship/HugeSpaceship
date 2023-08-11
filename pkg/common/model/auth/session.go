package auth

import (
	"HugeSpaceship/pkg/common/model/common"
	"net/netip"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID     uuid.UUID
	Username   string
	Game       common.GameType
	IP         netip.Addr // This is so that we can track changes in ip and force a re-auth, helps prevent token theft.
	Token      string
	ExpiryDate time.Time
	Platform   common.Platform
}
