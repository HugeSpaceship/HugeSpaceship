package lbp_xml

import "encoding/xml"

type NpData struct {
	XMLName xml.Name   `xml:"npdata"`
	Friends []NpHandle `xml:"friends"`
	Blocked []NpHandle `xml:"blocked"`
}
