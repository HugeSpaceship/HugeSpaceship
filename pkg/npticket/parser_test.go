package npticket

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	ticket, err := os.Open("../../test/example-ticket")
	if err != nil {
		t.Error(err)
	}

	parsedTicket, err := Parse(ticket)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(strings.Trim(string(parsedTicket.Username[:]), "\000"))
}
