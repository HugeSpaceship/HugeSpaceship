package npticket

import (
	"HugeSpaceship/testdata"
	"testing"
)

func TestVerifyTicket(t *testing.T) {
	parser := NewParser(testdata.ExamplePSNTicket)
	ticket, err := parser.Parse()
	if err != nil {
		t.Error(err)
	}

	if !VerifyTicket(ticket) {
		t.Error("Ticket is not valid")
	}

	parser2 := NewParser(testdata.ExampleRPCNTicket)
	ticket2, err := parser2.Parse()
	if err != nil {
		t.Error(err)
	}

	if !VerifyTicket(ticket2) {
		t.Error("Ticket is not valid")
	}

}
