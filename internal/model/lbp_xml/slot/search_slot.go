package slot

import (
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/model/lbp_xml/npdata"
	"encoding/xml"
	"github.com/google/uuid"
)

type SearchSlot struct {
	XMLName xml.Name `xml:"slot"`
	Type    string   `xml:"type,attr" hs_db:"-"`

	ID                  string          `xml:"id"`
	NPHandle            npdata.NpHandle `xml:"npHandle"`
	Uploader            uuid.UUID       `xml:"-" db:"uploader"`
	SearchScore         float32         `xml:"searchScore"`
	Location            common.Location `xml:"location"`
	LocationX           int32           `xml:"-" db:"location_x"`
	LocationY           int32           `xml:"-" db:"location_y"`
	Game                int             `xml:"game"`
	Name                string          `xml:"name,omitempty"`
	Description         string          `xml:"description,omitempty"`
	RootLevel           string          `xml:"rootLevel"`
	Icon                string          `xml:"icon"`
	InitiallyLocked     bool            `xml:"initiallyLocked"`
	IsSubLevel          bool            `xml:"isSubLevel" db:"sub_level"`
	IsLBP1Only          bool            `xml:"isLBP1Only" db:"lbp1only"`
	Background          string          `xml:"background"`
	Shareable           int             `xml:"shareable"`
	MinPlayers          uint            `xml:"minPlayers"`
	MaxPlayers          uint            `xml:"maxPlayers"`
	HeartCount          int64           `xml:"heartCount"`
	ThumbsUp            int64           `xml:"thumbsup" db:"thumbs_up_count"`
	ThumbsDown          int64           `xml:"thumbsdown" db:"thumbs_down_count"`
	AverageRating       float32         `xml:"averageRating"`
	PlayerCount         uint64          `xml:"playerCount"`
	MatchingPlayers     uint64          `xml:"matchingPlayers"`
	MMPick              bool            `xml:"mmpick"`
	IsAdventurePlanet   bool            `xml:"isAdventurePlanet"`
	Ps4Only             bool            `xml:"ps4Only"`
	PlayCount           uint64          `xml:"playCount" db:"total_play_count"`
	CompletionCount     uint64          `xml:"completionCount"`
	Lbp1PlayCount       uint64          `xml:"lbp1PlayCount"`
	Lbp1CompletionCount uint64          `xml:"lbp1CompletionCount"`
	Lbp1UniquePlayCount uint64          `xml:"lbp1UniquePlayCount"`
	Lbp2PlayCount       uint64          `xml:"lbp2PlayCount"`
	Lbp2CompletionCount uint64          `xml:"lbp2CompletionCount"`
	UniquePlayCount     uint64          `xml:"uniquePlayCount" db:"play_count"`
	Lbp3PlayCount       uint64          `xml:"lbp3PlayCount"`
	Lbp3CompletionCount uint64          `xml:"lbp3CompletionCount"`
	Lbp3UniquePlayCount uint64          `xml:"lbp3UniquePlayCount"`
}
