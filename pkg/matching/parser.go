package matching

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/matching/types"
	"github.com/rs/zerolog/log"
	"net/netip"
	"regexp"
	"strconv"
)

var JsonifyMatchRegex = regexp.MustCompile("^\\[([^,]*),\\[(.*)]]")
var HexLocationRegex = regexp.MustCompile("0x[a-fA-F0-9]{7,8}")

// hexToDec converts a hex utf-8 number (or technically any number) to a decimal utf-8 number
// this is used because JSON doesn't like hex number literals
func hexToDec(in []byte) []byte {
	parseInt, err := strconv.ParseUint(string(in), 0, 32)
	if err != nil {
		log.Debug().Bytes("in", in).Err(err).Msg("Failed to parse int")
	}
	return []byte(strconv.FormatUint(parseInt, 10))
}

// converts a match message into valid json
// uses a regex to pull out the command and the body
func jsonifyMatch(data []byte) []byte {
	// Will pull two things out of the match, the command name, and the body
	matches := JsonifyMatchRegex.FindSubmatch(data)

	buf := new(bytes.Buffer)
	buf.WriteRune('{')
	// this abomination takes an ip address that is represented by hex and turns it into a decimal number
	buf.Write(HexLocationRegex.ReplaceAllFunc(matches[2], hexToDec))
	buf.WriteRune('}')

	return buf.Bytes()
}

// Unmarshal parses a match message by first using regex to convert it to json
func Unmarshal[T any](data []byte) (T, error) {
	out := new(T)
	err := json.Unmarshal(jsonifyMatch(data), out)
	return *out, err
}

func GetCommand(data []byte) types.MatchCommand {
	matches := JsonifyMatchRegex.FindSubmatch(data)

	return types.MatchCommand(matches[1])
}

func IPFromLocation(location uint32) (netip.Addr, bool) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, location)
	if err != nil {
		return netip.Addr{}, false
	}

	return netip.AddrFromSlice(buf.Bytes())

}
