package types

type Ticket struct {
	MajorVersion uint8
	MinorVersion uint8

	Padding      [5]byte // Find out what this is, it may be important
	TicketLength uint8

	BodySectionHeader [4]byte

	SerielIDHeader [4]byte
	SerialID       [20]byte
	IssuerIDHeader [4]byte
	IssuerID       uint32

	IssuedDateHeader [4]byte
	IssuedDate       uint64
	ExpiryDateHeader [4]byte
	ExpiryDate       uint64

	UserIDHeader   [4]byte
	UserID         uint64
	UsernameHeader [4]byte
	Username       [20]byte

	Country [4]byte
	Domain  [4]byte

	TitleID [20]byte // This length is a guess

	Status uint32
}
