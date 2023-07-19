// Package npticket parses and verifies PSN NPTickets sent to us by sony.
// It supports RPCN and PSN v2.1 tickets for maximum compatibility
// PSN ticket v3 support will be added at a later date for PSVita compatibility.
package npticket

import (
	"HugeSpaceship/pkg/npticket/types"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"
	"time"
)

type TicketParser struct {
	reader     io.Reader
	ticketData []byte
	TicketBody []byte
	ticket     types.Ticket
}

func NewParser(data []byte) TicketParser {
	return TicketParser{reader: bytes.NewReader(data), ticketData: data}
}

func (parser TicketParser) ReadTicketHeader() (types.TicketHeader, error) {
	var header types.TicketHeader
	err := binary.Read(parser.reader, binary.BigEndian, &header)
	header.MajorVersion >>= 4
	return header, err
}

func (parser TicketParser) ReadDataHeader() (types.TicketDataHeader, error) {
	var header types.TicketDataHeader
	err := binary.Read(parser.reader, binary.BigEndian, &header)

	return header, err
}

func (parser TicketParser) ReadSectionHeader() (types.TicketSectionData, error) {
	header := types.TicketSectionData{}
	var headerData = make([]byte, 4)
	err := binary.Read(parser.reader, binary.BigEndian, &headerData)

	header.Type = types.SectionType(headerData[1])

	header.Length = headerData[3]

	return header, err
}

func (parser TicketParser) ReadFooter() (types.TicketFooter, error) {
	var header types.TicketFooter

	cipherID, err := parser.ReadBytes()
	if err != nil {
		return types.TicketFooter{}, err
	}

	signature, err := parser.ReadBytes()

	header.Signatory = binary.BigEndian.Uint32(cipherID)
	header.Signature = signature

	return header, err
}

func (parser TicketParser) ReadBytes() ([]byte, error) {
	header, err := parser.ReadDataHeader()
	if err != nil {
		return nil, err
	}
	output := make([]byte, header.Length)
	err = binary.Read(parser.reader, binary.BigEndian, &output)

	return output, err
}

func (parser TicketParser) ReadString() (string, error) {
	header, err := parser.ReadDataHeader()
	if err != nil {
		return "", err
	}
	output := make([]byte, header.Length)
	err = binary.Read(parser.reader, binary.BigEndian, &output)

	return strings.Trim(string(output), "\000"), err
}

func (parser TicketParser) ReadUInt32() (uint32, error) {
	header, err := parser.ReadDataHeader()
	if err != nil {
		return 0, err
	}
	if header.Type != types.UInt32 {
		return 0, errors.New("mismatched type expected uint32, got " + types.TypeToString(header.Type))
	}
	var output uint32
	err = binary.Read(parser.reader, binary.BigEndian, &output)

	return output, err
}

func (parser TicketParser) ReadUInt64() (uint64, error) {
	header, err := parser.ReadDataHeader()
	if err != nil {
		return 0, err
	}
	if header.Type != types.UInt64 {
		return 0, errors.New("mismatched type expected UInt64, got " + types.TypeToString(header.Type) + "(")
	}
	var output uint64
	err = binary.Read(parser.reader, binary.BigEndian, &output)
	return output, err
}

func (parser TicketParser) ReadTimestamp() (time.Time, error) {
	header, err := parser.ReadDataHeader()
	if err != nil {
		return time.Time{}, err
	}
	if header.Type != types.Timestamp {
		return time.Time{}, errors.New("mismatched type expected timestamp, got " + types.TypeToString(header.Type))
	}
	var output uint64
	err = binary.Read(parser.reader, binary.BigEndian, &output)
	return time.UnixMilli(int64(output)), err
}

func (parser TicketParser) Parse() (types.Ticket, error) {

	var ticket = types.Ticket{}

	header, err := parser.ReadTicketHeader()

	if err != nil {
		return types.Ticket{}, err
	}
	_, _ = parser.ReadSectionHeader()
	if header.MajorVersion == 2 && header.MinorVersion == 1 {
		ticket, err = parser.parseVersion2Ticket()
		if err != nil {
			return types.Ticket{}, err
		}
	} else if header.MajorVersion == 3 && header.MinorVersion == 0 {
		ticket, err = parser.parseVersion3Ticket()
		if err != nil {
			return types.Ticket{}, err
		}
	}

	parser.ticket = ticket
	return ticket, nil
}
