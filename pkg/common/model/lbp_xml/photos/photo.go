package photos

import (
	"encoding/xml"
)

/*
<photo timestamp="1507929818678">

	  <id>126062122</id>
	  <author>AntDRAG</author>
	  <small>2522288e272f704b422cc8c6ad2028920e6eac86</small>
	  <medium>2522288e272f704b422cc8c6ad2028920e6eac86</medium>
	  <large>c7b2083de1cf53a71cd4b7c9e9c6fd9e10959c8f</large>
	  <plan>cfaa379c7b2c82bd8cfda67bc894eca3ec1ecdc1</plan>
	  <slot type="developer">
	    <id>0</id>
	  </slot>
	</photo>
*/
type Photos struct {
	XMLName xml.Name `xml.Name:"photos"`
	Photos  []Photo  `xml:"slots,omitempty"`
}

type Photo struct {
	XMLName   xml.Name  `xml:"photo"`
	Timestamp int64     `xml:"timestamp,attr"`
	ID        uint64    `xml:"id"`
	Author    string    `xml:"author"`
	Small     string    `xml:"small"`
	Medium    string    `xml:"medium"`
	Large     string    `xml:"large"`
	Plan      string    `xml:"plan"`
	Slot      PhotoSlot `xml:"slot,omitempty"`
}
