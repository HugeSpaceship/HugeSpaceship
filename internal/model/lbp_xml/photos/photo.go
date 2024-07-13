package photos

import (
	"encoding/xml"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml/npdata"
	"strconv"
	"strings"
)

// Photos is an LBP XML friendly list of photos
type Photos struct {
	XMLName xml.Name `xml.Name:"photos"`
	Photos  []Photo  `xml:"slots,omitempty"`
}

// Photo is the structure for LBP photos
type Photo struct {
	XMLName   xml.Name      `xml:"photo"`
	Timestamp int64         `xml:"timestamp,attr"`
	ID        uint64        `xml:"id,omitempty"`
	Author    string        `xml:"author,omitempty"`
	Small     string        `xml:"small"`
	Medium    string        `xml:"medium"`
	Large     string        `xml:"large"`
	Plan      string        `xml:"plan"`
	Slot      PhotoSlot     `xml:"slot,omitempty"`
	Subjects  PhotoSubjects `xml:"subjects"`
}

// PhotoSubjects is an LBP XML friendly list of photo subjects
type PhotoSubjects struct {
	Subjects []PhotoSubject
}

// PhotoSubject is a player within a photo
type PhotoSubject struct {
	NpHandle    npdata.NpHandle `xml:"npHandle"`
	DisplayName string          `xml:"displayName"`
	Bounds      string          `xml:"bounds"`
}

// PhotoSubjectBounds encapsulates some utilities for parsing the comma seperated values in the xml
type PhotoSubjectBounds string

func (subject PhotoSubjectBounds) GetBounds() (x1, y1, x2, y2 float64, err error) {
	parts := strings.Split(string(subject), ",")
	x1, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return
	}
	y1, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return
	}
	x2, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return
	}
	y2, err = strconv.ParseFloat(parts[3], 64)
	return
}
