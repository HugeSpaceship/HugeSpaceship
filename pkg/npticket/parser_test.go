package npticket

import (
	"github.com/HugeSpaceship/HugeSpaceship/testdata"
	"testing"
)

func TestParse(t *testing.T) {
	parser := NewParser(testdata.Wack)

	ticket, err := parser.Parse()
	if err != nil {
		t.Error(err)
	}
	if ticket.Username != "Zaprit282" {
		t.Error("Username does not match expected value")
	}
}
