package npticket

import (
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	ticketData, err := os.ReadFile("../../test/example-ticket")
	if err != nil {
		t.Error(err)
	}

	parser := NewParser(ticketData)
	if err != nil {
		t.Error(err)
	}
	ticket, err := parser.Parse()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(ticket.Username)
}
