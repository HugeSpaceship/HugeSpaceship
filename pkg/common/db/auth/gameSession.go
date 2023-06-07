package auth

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model"
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/npticket/types"
	"crypto/sha256"
	"encoding/base64"
	"github.com/google/uuid"
	"net/netip"
	"time"
)

var CachedSessions = make(map[string]auth.Session)

func NewSession(ticket types.Ticket, ip netip.Addr) string {
	c := db.GetConnection()

	id := uuid.New().String()

	sum := sha256.Sum256([]byte(id))
	token := base64.URLEncoding.EncodeToString(sum[:])

	CachedSessions[token] = auth.Session{
		ID:         0,
		UserID:     0,
		Username:   ticket.Username,
		Game:       model.LBP2,
		IP:         ip,
		Token:      token,
		ExpiryDate: time.Now().Add(24 * time.Hour),
	}

	err := c.NewSession(ticket.Username, model.LBP2, ip, model.PS3, id)
	if err != nil {
		panic(err.Error())
		return token
	}

	return token
}

func GetSession(token string) auth.Session {
	var session auth.Session
	// TODO: return session
	return session
}

func IsSessionValid(token string) bool {
	_, exists := CachedSessions[token]
	return exists
}
