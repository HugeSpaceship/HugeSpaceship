package npticket

import (
	"os"
	"testing"
)

func TestVerifyTicket(t *testing.T) {
	data, err := os.ReadFile("../../test/example-ticket")
	if err != nil {
		t.Error(err)
	}
	parser := NewParser(data)
	_, err = parser.Parse()
	if err != nil {
		t.Error(err)
	}

	if !parser.VerifyTicket() {
		t.Error("Ticket is not valid")
	}

}
