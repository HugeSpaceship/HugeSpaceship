package news

import (
	"github.com/google/uuid"
	"time"
)

type News struct {
	ID         uuid.UUID `db:"id"`
	Title      string
	Content    string
	Summary    string
	Date       time.Time
	Image      string
	Category   uuid.UUID
	Background string
}
