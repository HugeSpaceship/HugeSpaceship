package lbp_xml

import (
	"HugeSpaceship/pkg/common/model/lbp_xml/slot"
	"encoding/xml"
)

type UploadPhoto struct {
	XMLName   xml.Name        `xml:"photo"`
	Timestamp int64           `xml:"timestamp,attr"`
	Small     string          `xml:"small"`
	Medium    string          `xml:"medium"`
	Large     string          `xml:"large"`
	Plan      string          `xml:"plan"`
	Slot      *slot.PhotoSlot `xml:"slot,omitempty"`
}
