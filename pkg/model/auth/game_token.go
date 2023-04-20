package auth

import "time"

type LoginResult struct {
	AuthTicket string `xml:"authTicket"`
	LBPEnvVer  string `xml:"lbpEnvVer"`
}

type GameSession struct {
	User uint64

	ExpiryDate time.Time
}
