package npdata

import (
	"database/sql"
	"encoding/xml"
)

// NpHandle is a field used in LBPs xml for storing the username and icon of a player
type NpHandle struct {
	XMLName    xml.Name         `xml:"npHandle"`
	Username   string           `xml:",innerxml"`
	IconHash   string           `xml:"icon,attr,omitempty" db:"-"`
	IconHashDB sql.Null[string] `db:"avatar_hash" xml:"-"`
}
