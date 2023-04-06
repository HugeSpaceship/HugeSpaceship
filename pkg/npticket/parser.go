package npticket

import (
	"HugeSpaceship/pkg/npticket/types"
	"encoding/binary"
	"io"
)

func Parse(reader io.Reader) (types.Ticket, error) {
	var ticket types.Ticket
	err := binary.Read(reader, binary.BigEndian, &ticket)
	if err != nil {
		return types.Ticket{}, err
	}

	return ticket, nil
}
