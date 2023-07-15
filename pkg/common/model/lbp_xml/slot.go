package lbp_xml

import (
	"encoding/xml"
	"github.com/google/uuid"
	"time"
)

type Slots struct {
	XMLName   xml.Name `xml:"slots"`
	Total     int      `xml:"total,attr"`
	HintStart int      `xml:"hint_start,attr"`
	Slots     []Slot
}

type SlotData struct {
	XMLName           xml.Name  `xml:"slot" db:"-"`
	ID                uuid.UUID `xml:"id,omitempty"`
	Type              string    `xml:"type,attr" db:"-"`
	Name              string    `xml:"name,omitempty" db:"name"`
	Description       string    `xml:"description,omitempty" db:"description"`
	Icon              string    `xml:"icon,omitempty" db:"icon"`
	RootLevel         string    `xml:"rootLevel,omitempty" db:"root_level"`
	Resource          []string  `xml:"resource" db:"-"`
	Location          Location  `xml:"location,omitempty" db:"-"`
	LocationX         int       `xml:"-" db:"locationX"`
	LocationY         int       `xml:"-" db:"locationY"`
	InitiallyLocked   bool      `xml:"initiallyLocked,omitempty" db:"initially_locked"`
	IsSubLevel        bool      `xml:"isSubLevel,omitempty" db:"sub_level"`
	IsLBP1Only        bool      `xml:"isLBP1Only,omitempty" db:"lbp1only"`
	Shareable         int       `xml:"shareable,omitempty" db:"shareable"`
	Background        string    `xml:"background,omitempty" db:"background"`
	Links             string    `xml:"links,omitempty" db:"-"`
	InternalLinks     string    `xml:"internalLinks,omitempty" db:"-"`
	LevelType         string    `xml:"levelType,omitempty" db:"level_type"`
	MinPlayers        int       `xml:"minPlayers,omitempty" db:"min_players"`
	MaxPlayers        int       `xml:"maxPlayers,omitempty" db:"max_players"`
	MoveRequired      bool      `xml:"moveRequired,omitempty" db:"move_required"`
	FirstPublished    time.Time `xml:"-" db:"first_published"`
	FirstPublishedXML int64     `xml:"first_published" db:"-"`
	LastUpdated       time.Time `xml:"-" db:"last_updated"`
	LastUpdatedXML    int64     `xml:"last_updated" db:"-"`
	Domain            int       `xml:"-" db:"domain"`
}

type StartPublishSlotResponse struct {
	XMLName  xml.Name `xml:"slot"`
	Resource []string `xml:"resource"`
}

type Slot struct {
	XMLName xml.Name `xml:"slot"`
	SlotData

	NpHandle          NpHandle `xml:"npHandle,omitempty"`
	Game              int      `xml:"game,omitempty"`
	HeartCount        int      `xml:"heartCount,omitempty"`
	ThumbsUp          int      `xml:"thumbsup,omitempty"`
	ThumbsDown        int      `xml:"thumbsdown,omitempty"`
	AverageRating     float32  `xml:"averageRating,omitempty"`
	PlayerCount       int      `xml:"playerCount,omitempty"`
	MatchingPlayers   int      `xml:"matchingPlayers,omitempty"`
	TeamPick          bool     `xml:"mmpick,omitempty"`
	CommentsEnabled   bool     `xml:"commentsEnabled,omitempty"`
	ReviewsEnabled    bool     `xml:"reviewsEnabled,omitempty"`
	UserLBP1PlayCount int      `xml:"yourlbp1PlayCount,omitempty"`
	UserLBP2PlayCount int      `xml:"yourlbp2PlayCount,omitempty"`
	PublishedIn       string   `xml:"publishedIn,omitempty"`

	ReviewCount      int `xml:"reviewCount,omitempty"`
	CommentCount     int `xml:"commentCount,omitempty"`
	PhotoCount       int `xml:"photoCount,omitempty"`
	AuthorPhotoCount int `xml:"authorPhotoCount,omitempty"`
	PlayCount        int `xml:"playCount,omitempty"`
	UniquePlayCount  int `xml:"uniquePlayCount,omitempty"`
	CompletionCount  int `xml:"completionCount,omitempty"`

	LBP1PlayCount       int `xml:"lbp1PlayCount,omitempty"`
	LBP1CompletionCount int `xml:"lbp1CompletionCount,omitempty"`
	LBP1UniquePlayCount int `xml:"lbp1UniquePlayCount,omitempty"`
	LBP2PlayCount       int `xml:"lbp2PlayCount,omitempty"`
	LBP2CompletionCount int `xml:"lbp2CompletionCount,omitempty"`
	LBP2UniquePlayCount int `xml:"lbp2UniquePlayCount,omitempty"`
	LBP3PlayCount       int `xml:"lbp3PlayCount,omitempty"`
	LBP3CompletionCount int `xml:"lbp3CompletionCount,omitempty"`
	LBP3UniquePlayCount int `xml:"lbp3UniquePlayCount,omitempty"`
}

// Location represents a place on a users earth... I think
type Location struct {
	X int `xml:"x"`
	Y int `xml:"y"`
}

type SlotResource struct {
	SlotID       uuid.UUID `db:"slot_id"`
	ResourceHash string    `db:"resource_hash"`
}
