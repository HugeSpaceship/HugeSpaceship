package types

import (
	"time"
)

type Ticket struct {
	SerialID string

	IssuerID uint32

	IssuedDate time.Time
	ExpiryDate time.Time

	UserID   uint64
	Username string

	Country string
	Domain  string

	TitleID string

	Status uint32

	Footer TicketFooter

	BodyHash []byte
}
