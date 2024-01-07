package slot

import (
	"HugeSpaceship/internal/model/common"
	"encoding/xml"
	"github.com/google/uuid"
	"time"
)

type Score struct {
	ID       uuid.UUID       `xml:"-" db:"id"`
	UserID   uuid.UUID       `xml:"-" db:"user_id"`
	Time     time.Time       `xml:"-" db:"achieved_time"`
	Platform common.Platform `xml:"-" db:"platform"`

	Type       int      `xml:"type" db:"type"`
	PlayerIds  []string `xml:"playerIds" db:"players"`
	MainPlayer string   `xml:"mainPlayer" db:"main_player"`
	Rank       uint64   `xml:"rank" db:"-"`
	Score      uint32   `xml:"score" db:"score"`
}

type ScoreBoard struct {
	XMLName     xml.Name `xml:"scores"`
	PlayRecord  []*Score `xml:"playRecord"`
	TotalScores uint64   `xml:"totalNumScores"`
	YourScore   uint32   `xml:"yourScore"`
	YourRank    uint64   `xml:"yourRank"`
}
