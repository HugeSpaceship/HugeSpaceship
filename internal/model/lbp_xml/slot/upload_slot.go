package slot

import (
	"encoding/xml"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
)

type Upload struct {
	XMLName         xml.Name        `xml:"slot" `
	ID              uint64          `xml:"id,omitempty"`
	Type            string          `xml:"type,attr" `
	Name            string          `xml:"name,omitempty"`
	Description     string          `xml:"description,omitempty"`
	Icon            string          `xml:"icon,omitempty"`
	RootLevel       string          `xml:"rootLevel,omitempty"`
	Resources       []string        `xml:"resource"`
	Location        common.Location `xml:"location,omitempty"`
	InitiallyLocked bool            `xml:"initiallyLocked,omitempty"`
	IsSubLevel      bool            `xml:"isSubLevel,omitempty"`
	IsLBP1Only      bool            `xml:"isLBP1Only,omitempty"`
	Shareable       int             `xml:"shareable,omitempty"`
	Background      string          `xml:"background,omitempty"`
	Links           string          `xml:"links,omitempty"`
	InternalLinks   string          `xml:"internalLinks,omitempty"`
	LevelType       string          `xml:"levelType,omitempty"`
	MinPlayers      int             `xml:"minPlayers,omitempty"`
	MaxPlayers      int             `xml:"maxPlayers,omitempty"`
	MoveRequired    bool            `xml:"moveRequired,omitempty"`
}
