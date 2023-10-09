package lbp_xml

import "encoding/xml"

type Notification struct {
	XMLName  xml.Name
	Type     string `xml:"type,attr"`
	Text     string `xml:"text"`
	Extended string `xml:"extended"`
}
