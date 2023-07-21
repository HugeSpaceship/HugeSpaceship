package slot

import (
	"HugeSpaceship/pkg/common/model/common"
	"HugeSpaceship/pkg/common/model/lbp_xml/npdata"
	"encoding/xml"
	"time"
)

type Type interface {
	Slot | Upload | SearchSlot
}

type Slots[S Type] struct {
	XMLName   xml.Name `xml:"slots"`
	Total     int      `xml:"total,attr"`
	HintStart int      `xml:"hint_start,attr"`
	Slots     []any
}

type StartPublishSlotResponse struct {
	XMLName  xml.Name `xml:"slot"`
	Resource []string `xml:"resource"`
}

type Slot struct {
	XMLName             xml.Name        `xml:"slot"`
	Type                string          `xml:"type,attr"`
	ID                  uint64          `xml:"id"`
	NpHandle            npdata.NpHandle `xml:"npHandle"`
	Location            common.Location `xml:"location"`
	Game                uint            `xml:"game"`
	RootLevel           string          `xml:"rootLevel"`
	Resources           []string        `xml:"resource"`
	Icon                string          `xml:"icon"`
	InitiallyLocked     bool            `xml:"initiallyLocked"`
	IsSubLevel          bool            `xml:"isSubLevel"`
	IsLBP1Only          bool            `xml:"isLBP1Only"`
	Background          string          `xml:"background"`
	Shareable           uint            `xml:"shareable"`
	MinPlayers          uint            `xml:"minPlayers"`
	MaxPlayers          uint            `xml:"maxPlayers"`
	HeartCount          uint64          `xml:"heartCount"`
	ThumbsUp            uint64          `xml:"thumbsup"`
	ThumbsDown          uint64          `xml:"thumbsdown"`
	AverageRating       float32         `xml:"averageRating"`
	PlayerCount         uint64          `xml:"playerCount"`
	MatchingPlayers     uint64          `xml:"matchingPlayers"`
	MMPick              bool            `xml:"mmpick"`
	YourRating          int             `xml:"yourRating"`
	YourDPadRating      int             `xml:"yourDPadRating"`
	YourLBP1PlayCount   uint64          `xml:"yourlbp1PlayCount"`
	YourLBP2PlayCount   uint64          `xml:"yourlbp2PlayCount"`
	ReviewCount         uint64          `xml:"reviewCount"`
	CommentCount        uint64          `xml:"commentCount"`
	PhotoCount          uint64          `xml:"photoCount"`
	AuthorPhotoCount    uint64          `xml:"authorPhotoCount"`
	FirstPublishedXML   int64           `xml:"firstPublished"`
	LastUpdatedXML      int64           `xml:"lastUpdated"`
	FirstPublished      time.Time       `xml:"-"`
	LastUpdated         time.Time       `xml:"-"`
	BadgeURL            string          `xml:"badgeURL"`
	CommentsEnabled     bool            `xml:"commentsEnabled"`
	ReviewsEnabled      bool            `xml:"reviewsEnabled"`
	PublishedIn         common.GameType `xml:"publishedIn"`
	PlayCount           uint64          `xml:"playCount"`
	CompletionCount     uint64          `xml:"completionCount"`
	Lbp1PlayCount       uint64          `xml:"lbp1PlayCount"`
	Lbp1CompletionCount uint64          `xml:"lbp1CompletionCount"`
	Lbp1UniquePlayCount uint64          `xml:"lbp1UniquePlayCount"`
	Lbp2PlayCount       uint64          `xml:"lbp2PlayCount"`
	Lbp2CompletionCount uint64          `xml:"lbp2CompletionCount"`
	UniquePlayCount     uint64          `xml:"uniquePlayCount"`
	Lbp3PlayCount       uint64          `xml:"lbp3PlayCount"`
	Lbp3CompletionCount uint64          `xml:"lbp3CompletionCount"`
	Lbp3UniquePlayCount uint64          `xml:"lbp3UniquePlayCount"`
}

// Location represents a place on a users earth... I think

type SlotResource struct {
	SlotID       int    `db:"slot_id"`
	ResourceHash string `db:"resource_hash"`
}
