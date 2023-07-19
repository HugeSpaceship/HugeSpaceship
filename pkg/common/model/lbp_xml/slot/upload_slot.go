package slot

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"encoding/xml"
)

type Upload struct {
	XMLName         xml.Name         `xml:"slot" `
	Type            string           `xml:"type,attr" `
	Name            string           `xml:"name,omitempty"`
	Description     string           `xml:"description,omitempty"`
	Icon            string           `xml:"icon,omitempty"`
	RootLevel       string           `xml:"rootLevel,omitempty"`
	Resource        []string         `xml:"resource"`
	Location        lbp_xml.Location `xml:"location,omitempty"`
	InitiallyLocked bool             `xml:"initiallyLocked,omitempty"`
	IsSubLevel      bool             `xml:"isSubLevel,omitempty"`
	IsLBP1Only      bool             `xml:"isLBP1Only,omitempty"`
	Shareable       int              `xml:"shareable,omitempty"`
	Background      string           `xml:"background,omitempty"`
	Links           string           `xml:"links,omitempty"`
	InternalLinks   string           `xml:"internalLinks,omitempty"`
	LevelType       string           `xml:"levelType,omitempty"`
	MinPlayers      int              `xml:"minPlayers,omitempty"`
	MaxPlayers      int              `xml:"maxPlayers,omitempty"`
	MoveRequired    bool             `xml:"moveRequired,omitempty"`
}
