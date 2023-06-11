package npticket

import (
	"HugeSpaceship/pkg/npticket/types"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

func (parser TicketParser) parseVersion3Ticket() (types.Ticket, error) {
	ticket := types.Ticket{}
	footer := types.TicketFooter{}
	var err error

	ticket.SerialID, err = parser.ReadString()
	if err != nil {
		return types.Ticket{}, err
	}
	ticket.IssuerID, err = parser.ReadUInt32()
	if err != nil {
		return types.Ticket{}, err
	}
	ticket.IssuedDate, err = parser.ReadTimestamp()
	if err != nil {
		return types.Ticket{}, err
	}
	ticket.ExpiryDate, err = parser.ReadTimestamp()
	if err != nil {
		return types.Ticket{}, err
	}
	ticket.UserID, err = parser.ReadUInt64()
	if err != nil {
		return types.Ticket{}, err
	}
	ticket.Username, err = parser.ReadString()
	if err != nil {
		return types.Ticket{}, err
	}
	ticket.Country, err = parser.ReadString()
	if err != nil {
		return types.Ticket{}, err
	}
	ticket.Domain, err = parser.ReadString()
	if err != nil {
		return types.Ticket{}, err
	}
	ticket.TitleID, err = parser.ReadString()
	if err != nil {
		return types.Ticket{}, err
	}

	_, _ = parser.ReadSectionHeader() // We don't care about date of birth
	_, _ = parser.ReadSectionHeader()
	_, _ = parser.ReadSectionHeader()
	_, _ = parser.ReadSectionHeader()
	_, _ = parser.ReadSectionHeader()
	_, _ = parser.ReadSectionHeader()

	footerHeader, err := parser.ReadSectionHeader()
	if err != nil {
		return types.Ticket{}, err
	}

	signatory, err := parser.ReadBytes()
	footer.Signatory = binary.BigEndian.Uint32(signatory)
	footer.Signature, err = parser.ReadBytes()
	ticket.Footer = footer

	switch footer.Signatory {
	case types.PSNSignatoryID:
		parser.TicketBody = parser.ticketData[:len(parser.ticketData)-0x38]
		digest := sha1.Sum(parser.TicketBody)
		ticket.BodyHash = digest[:]
	case types.RPCNSignatoryID:
		// the -4 at the end is because the header exists
		parser.TicketBody = parser.ticketData[8 : len(parser.ticketData)-int(footerHeader.Length)-4]
		digest := sha256.Sum224(parser.TicketBody)
		ticket.BodyHash = digest[:]
	default:
		return types.Ticket{}, errors.New("invalid signatory")
	}

	return ticket, nil
}
