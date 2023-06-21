package lbp_xml

import "encoding/xml"

type Slots struct {
	Total     int `xml:"total,attr"`
	HintStart int `xml:"hint_start,attr"`
	Slots     []Slot
}

type Slot struct {
	XMLName         xml.Name  `xml:"slot"`
	Type            string    `xml:"type,attr"`
	Name            string    `xml:"name,omitempty"`
	Description     string    `xml:"description,omitempty"`
	Icon            string    `xml:"icon,omitempty"`
	RootLevel       string    `xml:"rootLevel,omitempty"`
	Resource        []string  `xml:"resource"`
	Location        *Location `xml:"location,omitempty"`
	InitiallyLocked *bool     `xml:"initiallyLocked,omitempty"`
	IsSubLevel      *bool     `xml:"isSubLevel,omitempty"`
	IsLBP1Only      *bool     `xml:"isLBP1Only,omitempty"`
	Shareable       *int      `xml:"shareable,omitempty"`
	Background      string    `xml:"background,omitempty"`
	Links           string    `xml:"links,omitempty"`
	InternalLinks   string    `xml:"internalLinks,omitempty"`
	LevelType       string    `xml:"levelType,omitempty"`
	MinPlayers      *int      `xml:"minPlayers,omitempty"`
	MaxPlayers      *int      `xml:"maxPlayers,omitempty"`
	MoveRequired    *bool     `xml:"moveRequired,omitempty"`
}

// Location represents a place on a users earth... I think
type Location struct {
	X int `xml:"x"`
	Y int `xml:"y"`
}
