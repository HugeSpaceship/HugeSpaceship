package auth

import (
	"HugeSpaceship/pkg/common/model/auth"
	"time"
)

var CachedSessions map[string]auth.GameSession

func NewSession() {

}

func GetSession(token string) auth.GameSession { // TODO: Refactor this with what gets merged from my laptop
	// TODO: return session
	return auth.GameSession{ // shitty stub
		User:       0,
		Token:      token,
		ExpiryDate: time.Now().Add(24 * time.Hour),
	}
}

func IsSessionValid(token string) bool {
	return GetSession(token).Token == token
}
