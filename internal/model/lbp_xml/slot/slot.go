package slot

import (
	"HugeSpaceship/internal/model/common"
	"HugeSpaceship/internal/model/lbp_xml/npdata"
	"encoding/xml"
	"github.com/google/uuid"
	"time"
)

type Type interface {
	Slot | Upload | SearchSlot | UpdateUserSlot
}

type PaginatedSlotList[S Type] struct {
	XMLName   xml.Name `xml:"slots"`
	Total     uint64   `xml:"total,attr"`
	HintStart uint64   `xml:"hint_start,attr"`
	Slots     []S
}

type List[S Type] struct {
	XMLName xml.Name `xml:"slots"`
	Slots   []S
}

type StartPublishSlotResponse struct {
	XMLName  xml.Name `xml:"slot"`
	Resource []string `xml:"resource"`
}

type Slot struct {
	XMLName     xml.Name         `xml:"slot"`
	Type        string           `xml:"type,attr"`
	ID          uint64           `xml:"id"`
	UploaderID  uuid.UUID        `xml:"-" db:"uploader"`
	NpHandle    *npdata.NpHandle `xml:"npHandle"`
	Location    common.Location  `xml:"location"`
	LocationX   int32            `xml:"-" db:"location_x"`
	LocationY   int32            `xml:"-" db:"location_y"`
	Game        uint             `xml:"game"`
	Name        string           `xml:"name,omitempty"`
	Description string           `xml:"description,omitempty"`
	RootLevel   string           `xml:"rootLevel"`
	//Resources         []string        `xml:"resource"`
	Icon              string  `xml:"icon"`
	InitiallyLocked   bool    `xml:"initiallyLocked"`
	IsSubLevel        bool    `xml:"isSubLevel" db:"sub_level"`
	IsLBP1Only        bool    `xml:"isLBP1Only" db:"lbp1only"`
	Background        string  `xml:"background"`
	Shareable         uint    `xml:"shareable"`
	MinPlayers        uint    `xml:"minPlayers"`
	MaxPlayers        uint    `xml:"maxPlayers"`
	HeartCount        uint64  `xml:"heartCount"`
	ThumbsUp          uint64  `xml:"thumbsup" db:"thumbs_up_count"`
	ThumbsDown        uint64  `xml:"thumbsdown" db:"thumbs_down_count"`
	AverageRating     float32 `xml:"averageRating"`
	PlayerCount       uint64  `xml:"playerCount"`
	MatchingPlayers   uint64  `xml:"matchingPlayers"`
	MMPick            bool    `xml:"mmpick"`
	YourRating        int     `xml:"yourRating"`
	YourDPadRating    int     `xml:"yourDPadRating"`
	YourLBP1PlayCount uint64  `xml:"yourlbp1PlayCount"`
	YourLBP2PlayCount uint64  `xml:"yourlbp2PlayCount"`
	ReviewCount       uint64  `xml:"reviewCount"`
	CommentCount      uint64  `xml:"commentCount"`
	PhotoCount        uint64  `xml:"photoCount"`
	AuthorPhotoCount  uint64  `xml:"authorPhotoCount"`

	FirstPublishedXML string `xml:"firstPublished"` // These two are strings because otherwise it integer overflows
	LastUpdatedXML    string `xml:"lastUpdated"`

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

type SlotResource struct {
	SlotID       int    `hs_db:"slot_id"`
	ResourceHash string `hs_db:"resource_hash"`
}

type UpdateUserSlot struct {
	Type     string          `xml:"type,attr"`
	Id       uint64          `xml:"id"`
	Location common.Location `xml:"location"`
}
