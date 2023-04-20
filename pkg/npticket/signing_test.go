package npticket

import (
	"os"
	"testing"
)

func TestVerifyTicket(t *testing.T) {
	data, err := os.ReadFile("../../test/example-psn-ticket")
	if err != nil {
		t.Error(err)
	}
	parser := NewParser(data)
	ticket, err := parser.Parse()
	if err != nil {
		t.Error(err)
	}

	if !VerifyTicket(ticket) {
		t.Error("Ticket is not valid")
	}

	data2, err := os.ReadFile("../../test/example-rpcn-ticket")
	if err != nil {
		t.Error(err)
	}
	parser2 := NewParser(data2)
	ticket2, err := parser2.Parse()
	if err != nil {
		t.Error(err)
	}

	if !VerifyTicket(ticket2) {
		t.Error("Ticket is not valid")
	}

}
