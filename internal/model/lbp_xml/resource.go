package lbp_xml

import "encoding/xml"

type Resources struct {
	XMLName   xml.Name `xml:"resources"`
	Resources []string `xml:"resource"`
}
