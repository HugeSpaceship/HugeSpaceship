package slot

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"HugeSpaceship/pkg/common/model/lbp_xml/npdata"
	"encoding/xml"
)

type SearchSlot struct {
	XMLName xml.Name `xml:"slot"`
	Type    string   `xml:"type,attr" db:"-"`

	ID                  string           `xml:"id"`
	NPHandle            npdata.NpHandle  `xml:"npHandle"`
	SearchScore         float32          `xml:"searchScore"`
	Location            lbp_xml.Location `xml:"location"`
	Game                int              `xml:"game"`
	Name                string           `xml:"name"`
	Description         string           `xml:"description"`
	RootLevel           string           `xml:"rootLevel"`
	Icon                string           `xml:"icon"`
	InitiallyLocked     bool             `xml:"initiallyLocked"`
	IsSubLevel          bool             `xml:"isSubLevel"`
	IsLBP1Only          bool             `xml:"isLBP1Only"`
	Shareable           int              `xml:"shareable"`
	MinPlayers          uint             `xml:"minPlayers"`
	MaxPlayers          uint             `xml:"maxPlayers"`
	HeartCount          int64            `xml:"heartCount"`
	Thumbsup            int64            `xml:"thumbsup"`
	Thumbsdown          int64            `xml:"thumbsdown"`
	AverageRating       float32          `xml:"averageRating"`
	PlayerCount         uint64           `xml:"playerCount"`
	MatchingPlayers     uint64           `xml:"matchingPlayers"`
	MMPick              bool             `xml:"mmpick"`
	IsAdventurePlanet   bool             `xml:"isAdventurePlanet"`
	Ps4Only             bool             `xml:"ps4Only"`
	PlayCount           uint64           `xml:"playCount"`
	CompletionCount     uint64           `xml:"completionCount"`
	Lbp1PlayCount       uint64           `xml:"lbp1PlayCount"`
	Lbp1CompletionCount uint64           `xml:"lbp1CompletionCount"`
	Lbp1UniquePlayCount uint64           `xml:"lbp1UniquePlayCount"`
	Lbp2PlayCount       uint64           `xml:"lbp2PlayCount"`
	Lbp2CompletionCount uint64           `xml:"lbp2CompletionCount"`
	UniquePlayCount     uint64           `xml:"uniquePlayCount"`
	Lbp3PlayCount       uint64           `xml:"lbp3PlayCount"`
	Lbp3CompletionCount uint64           `xml:"lbp3CompletionCount"`
	Lbp3UniquePlayCount uint64           `xml:"lbp3UniquePlayCount"`
}
