package npticket

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	ticketData, err := os.ReadFile("../../test/example-psn-ticket")
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

	if ticket.Username != "Zaprit282" {
		t.Error("Username does not match expected value")
	}
}
