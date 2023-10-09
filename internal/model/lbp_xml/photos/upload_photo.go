package photos

import (
	"encoding/xml"
)

type UploadPhoto struct {
	XMLName   xml.Name      `xml:"photo"`
	Timestamp int64         `xml:"timestamp,attr"`
	Small     string        `xml:"small"`
	Medium    string        `xml:"medium"`
	Large     string        `xml:"large"`
	Plan      string        `xml:"plan"`
	Slot      PhotoSlot     `xml:"slot,omitempty"`
	Subjects  PhotoSubjects `xml:"subjects"`
}
