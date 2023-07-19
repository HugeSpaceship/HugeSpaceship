package slot

import (
	"encoding/xml"
)

type Slots struct {
	XMLName   xml.Name `xml:"slots"`
	Total     int      `xml:"total,attr"`
	HintStart int      `xml:"hint_start,attr"`
	Slots     []Slot
}

type StartPublishSlotResponse struct {
	XMLName  xml.Name `xml:"slot"`
	Resource []string `xml:"resource"`
}

type Slot struct {
	XMLName xml.Name `xml:"slot"`
	Upload

	HeartCount        int     `xml:"heartCount"`
	ThumbsUp          int     `xml:"thumbsup"`
	ThumbsDown        int     `xml:"thumbsdown"`
	AverageRating     float32 `xml:"averageRating"`
	PlayerCount       int     `xml:"playerCount"`
	MatchingPlayers   int     `xml:"matchingPlayers"`
	TeamPick          bool    `xml:"mmpick"`
	CommentsEnabled   bool    `xml:"commentsEnabled"`
	ReviewsEnabled    bool    `xml:"reviewsEnabled"`
	UserLBP1PlayCount int     `xml:"yourlbp1PlayCount"`
	UserLBP2PlayCount int     `xml:"yourlbp2PlayCount"`
	PublishedIn       string  `xml:"publishedIn"`

	ReviewCount      int `xml:"reviewCount"`
	CommentCount     int `xml:"commentCount"`
	PhotoCount       int `xml:"photoCount"`
	AuthorPhotoCount int `xml:"authorPhotoCount"`
	PlayCount        int `xml:"playCount"`
	UniquePlayCount  int `xml:"uniquePlayCount"`
	CompletionCount  int `xml:"completionCount"`

	LBP1PlayCount       int `xml:"lbp1PlayCount"`
	LBP1CompletionCount int `xml:"lbp1CompletionCount"`
	LBP1UniquePlayCount int `xml:"lbp1UniquePlayCount"`
	LBP2PlayCount       int `xml:"lbp2PlayCount"`
	LBP2CompletionCount int `xml:"lbp2CompletionCount"`
	LBP2UniquePlayCount int `xml:"lbp2UniquePlayCount"`
	LBP3PlayCount       int `xml:"lbp3PlayCount"`
	LBP3CompletionCount int `xml:"lbp3CompletionCount"`
	LBP3UniquePlayCount int `xml:"lbp3UniquePlayCount"`
}

// Location represents a place on a users earth... I think

type SlotResource struct {
	SlotID       int    `db:"slot_id"`
	ResourceHash string `db:"resource_hash"`
}
