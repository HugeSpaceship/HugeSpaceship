package npdata

import "encoding/xml"

// NpHandle is a field used in LBPs xml for storing the username and icon of a player
type NpHandle struct {
	XMLName  xml.Name `xml:"npHandle"`
	Username string   `xml:",innerxml"`
	IconHash string   `xml:"icon,attr,omitempty" db:"avatar_hash"`
}
