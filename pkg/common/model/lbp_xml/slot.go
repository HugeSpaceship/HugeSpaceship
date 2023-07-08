package lbp_xml

import (
	"encoding/xml"
	"github.com/google/uuid"
)

type Slots struct {
	Total     int `xml:"total,attr"`
	HintStart int `xml:"hint_start,attr"`
	Slots     []Slot
}

type SlotData struct {
	XMLName         xml.Name  `xml:"slot"`
	Type            string    `xml:"type,attr"`
	Name            string    `xml:"name,omitempty"`
	Description     string    `xml:"description,omitempty"`
	Icon            string    `xml:"icon,omitempty"`
	RootLevel       string    `xml:"rootLevel,omitempty"`
	Resource        []string  `xml:"resource"`
	Location        *Location `xml:"location,omitempty"`
	InitiallyLocked *bool     `xml:"initiallyLocked,omitempty"`
	IsSubLevel      *bool     `xml:"isSubLevel,omitempty"`
	IsLBP1Only      *bool     `xml:"isLBP1Only,omitempty"`
	Shareable       *int      `xml:"shareable,omitempty"`
	Background      string    `xml:"background,omitempty"`
	Links           string    `xml:"links,omitempty"`
	InternalLinks   string    `xml:"internalLinks,omitempty"`
	LevelType       string    `xml:"levelType,omitempty"`
	MinPlayers      *int      `xml:"minPlayers,omitempty"`
	MaxPlayers      *int      `xml:"maxPlayers,omitempty"`
	MoveRequired    *bool     `xml:"moveRequired,omitempty"`
}

type Slot struct {
	XMLName xml.Name `xml:"slot"`
	SlotData

	ID                uuid.UUID `xml:"id,omitempty"`
	NpHandle          *NpHandle `xml:"npHandle,omitempty"`
	Game              *int      `xml:"game,omitempty"`
	HeartCount        *int      `xml:"heartCount,omitempty"`
	ThumbsUp          *int      `xml:"thumbsup,omitempty"`
	ThumbsDown        *int      `xml:"thumbsdown,omitempty"`
	AverageRating     *float32  `xml:"averageRating,omitempty"`
	PlayerCount       *int      `xml:"playerCount,omitempty"`
	MatchingPlayers   *int      `xml:"matchingPlayers,omitempty"`
	TeamPick          *bool     `xml:"mmpick,omitempty"`
	CommentsEnabled   *bool     `xml:"commentsEnabled,omitempty"`
	ReviewsEnabled    *bool     `xml:"reviewsEnabled,omitempty"`
	UserLBP1PlayCount *int      `xml:"yourlbp1PlayCount,omitempty"`
	UserLBP2PlayCount *int      `xml:"yourlbp2PlayCount,omitempty"`
	PublishedIn       *string   `xml:"publishedIn,omitempty"`

	FirstPublished *int64 `xml:"firstPublished,omitempty"`
	LastUpdated    *int64 `xml:"lastUpdated,omitempty"`

	ReviewCount      *int `xml:"reviewCount,omitempty"`
	CommentCount     *int `xml:"commentCount,omitempty"`
	PhotoCount       *int `xml:"photoCount,omitempty"`
	AuthorPhotoCount *int `xml:"authorPhotoCount,omitempty"`
	PlayCount        *int `xml:"playCount,omitempty"`
	UniquePlayCount  *int `xml:"uniquePlayCount,omitempty"`
	CompletionCount  *int `xml:"completionCount,omitempty"`

	LBP1PlayCount       *int `xml:"lbp1PlayCount,omitempty"`
	LBP1CompletionCount *int `xml:"lbp1CompletionCount,omitempty"`
	LBP1UniquePlayCount *int `xml:"lbp1UniquePlayCount,omitempty"`
	LBP2PlayCount       *int `xml:"lbp2PlayCount,omitempty"`
	LBP2CompletionCount *int `xml:"lbp2CompletionCount,omitempty"`
	LBP2UniquePlayCount *int `xml:"lbp2UniquePlayCount,omitempty"`
	LBP3PlayCount       *int `xml:"lbp3PlayCount,omitempty"`
	LBP3CompletionCount *int `xml:"lbp3CompletionCount,omitempty"`
	LBP3UniquePlayCount *int `xml:"lbp3UniquePlayCount,omitempty"`
}

// Location represents a place on a users earth... I think
type Location struct {
	X int `xml:"x"`
	Y int `xml:"y"`
}
