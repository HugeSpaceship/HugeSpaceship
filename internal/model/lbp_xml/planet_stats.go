package lbp_xml

import "encoding/xml"

type PlanetStats struct {
	XMLName xml.Name `xml:"planetStats"`

	SlotsCount     uint64 `xml:"totalSlotCount"`
	TeamPicksCount uint64 `xml:"mmPicksCount"`
}
