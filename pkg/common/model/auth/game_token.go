package auth

import "time"

type LoginResult struct {
	AuthTicket string `xml:"authTicket"`
	LBPEnvVer  string `xml:"lbpEnvVer"`
}

type GameSession struct {
	User uint64

	Token string // Should be uuid

	ExpiryDate time.Time
}
