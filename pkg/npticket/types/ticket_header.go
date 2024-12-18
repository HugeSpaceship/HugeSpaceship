package types

import "fmt"

type TicketDataHeader struct {
	Type   DataType
	Length uint16
}

type TicketSectionData struct {
	Type   SectionType
	Length uint8
}

type TicketHeader struct {
	MajorVersion uint8
	MinorVersion uint8

	Padding      [5]byte // Find out what this is, it may be important
	TicketLength uint8
}

func (h TicketHeader) GetVersion() string {
	return fmt.Sprintf("%d.%d", h.MajorVersion, h.MinorVersion)
}
