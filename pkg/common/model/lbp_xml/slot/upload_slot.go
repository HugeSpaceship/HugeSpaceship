package slot

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"encoding/xml"
	"time"

	"github.com/google/uuid"
)

type Upload struct {
	XMLName     xml.Name         `xml:"slot" db:"-"`
	Type        string           `xml:"type,attr" db:"-"`
	ID          *int             `xml:"id,omitempty" db:"id"`
	NpHandle    lbp_xml.NpHandle `xml:"npHandle,omitempty" db:"-"`
	Location    lbp_xml.Location `xml:"location,omitempty" db:"-"`
	Uploader    uuid.UUID        `xml:"-" db:"uploader"`
	Game        int              `xml:"game,omitempty"`
	Name        string           `xml:"name,omitempty" db:"name"`
	Description string           `xml:"description,omitempty" db:"description"`

	Icon      string   `xml:"icon,omitempty" db:"icon"`
	RootLevel string   `xml:"rootLevel,omitempty" db:"root_level"`
	Resource  []string `xml:"resource" db:"-"`

	LocationX         int       `xml:"-" db:"locationx"`
	LocationY         int       `xml:"-" db:"locationy"`
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
	FirstPublishedXML int64     `xml:"firstPublished" db:"-"`
	LastUpdated       time.Time `xml:"-" db:"last_updated"`
	LastUpdatedXML    int64     `xml:"lastUpdated" db:"-"`
	Domain            int       `xml:"-" db:"domain"`
}
