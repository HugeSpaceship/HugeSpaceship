package recent_activity

import (
	"encoding/xml"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml/slot"
)

type Stream struct {
	XMLName        xml.Name       `xml:"stream"`
	StartTimestamp int64          `xml:"start_timestamp"`
	EndTimestamp   int64          `xml:"end_timestamp"`
	Groups         Groups         `xml:"groups"`
	Users          []lbp_xml.User `xml:"users,omitempty"`
	Slots          []slot.Slot    `xml:"slots,omitempty"`
	News           News           `xml:"news,omitempty"`
}

type Groups struct {
	Groups []Group
}

type Group struct {
	XMLName   xml.Name    `xml:"group"`
	Type      string      `xml:"type,attr"`
	NewsID    string      `xml:"news_id,omitempty"`
	Timestamp int64       `xml:"timestamp"`
	Events    GroupEvents `xml:"events,omitempty"`
	Subgroups Groups      `xml:"subgroups,omitempty"`
}

type GroupEvents struct {
	Events []Event
}

type Event struct {
	XMLName   xml.Name `xml:"event"`
	Type      string   `xml:"type,attr"`
	Timestamp int64    `xml:"timestamp"`
	NewsID    string   `xml:"news_id,omitempty"`
}

type News struct {
	Items []NewsItem
}

type NewsItem struct {
	XMLName    xml.Name
	ID         int    `xml:"id"`
	Title      string `xml:"title"`
	Text       string `xml:"text"`
	Summary    string `xml:"summary,omitempty"`
	Date       int64  `xml:"date"`
	Image      Image  `xml:"image"`
	Category   string `xml:"category"`
	Background string `xml:"background,omitempty"`
}

type Image struct {
	Alignment string `xml:"alignment"`
}
