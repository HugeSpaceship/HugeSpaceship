package photos

import "encoding/xml"

type PhotoSlot struct {
	XMLName     xml.Name `xml:"slot"`
	Type        string   `xml:"type,attr" db:"-"`
	ID          int64    `xml:"id"`
	Name        string   `xml:"name,omitempty"`
	Description string   `xml:"description,omitempty"`
}
