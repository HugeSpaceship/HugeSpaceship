package lbp_xml

import (
	"database/sql"
	"encoding/xml"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml/npdata"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml/slot"
	"github.com/google/uuid"
)

type User struct {
	XMLName    xml.Name        `db:"-"`
	ID         uuid.UUID       `xml:"-" db:"id"`
	PSN_UID    string          `xml:"-" db:"psn_uid"`
	RPCN_UID   string          `xml:"-" db:"rpcn_uid"`
	Type       string          `xml:"type,attr" db:"-"`
	NpHandle   npdata.NpHandle `xml:"npHandle" db:"-"`
	Username   string          `xml:"-" db:"username"`
	AvatarHash sql.NullString  `xml:"-" db:"avatar_hash"`
	Biography  string          `xml:"biography" db:"bio"`

	Planets    string `xml:"planets" db:"-"` // The resource of the planet for the correct game
	LBP2Planet string `xml:"-" db:"planet_lbp2"`
	LBP3Planet string `xml:"-" db:"planet_lbp3"`
	LBPVPlanet string `xml:"-" db:"planet_lbp_vita"`
	CCPlanet   string `xml:"-" db:"planet_cc"`

	Game                      string          `xml:"game"`
	Lbp1UsedSlots             int             `xml:"lbp1UsedSlots"`
	EntitledSlots             int             `xml:"entitledSlots" db:"entitled_slots"`
	FreeSlots                 int             `xml:"freeSlots" db:"free_slots"`
	CrossControlUsedSlots     int             `xml:"crossControlUsedSlots"`
	CrossControlEntitledSlots int             `xml:"crossControlEntitledSlots"`
	CrossControlFreeSlots     int             `xml:"crossControlFreeSlots"`
	Lbp2UsedSlots             int             `xml:"lbp2UsedSlots" db:"used_slots"`
	Lbp2EntitledSlots         int             `xml:"lbp2EntitledSlots"`
	Lbp2FreeSlots             int             `xml:"lbp2FreeSlots"`
	Lbp3UsedSlots             int             `xml:"lbp3UsedSlots"`
	Lbp3EntitledSlots         int             `xml:"lbp3EntitledSlots"`
	Lbp3FreeSlots             int             `xml:"lbp3FreeSlots"`
	Lists                     int             `xml:"lists"`
	ListsQuota                int             `xml:"lists_quota"`
	HeartCount                int             `xml:"heartCount"`
	ReviewCount               int             `xml:"reviewCount"`
	CommentCount              int             `xml:"commentCount"`
	PhotosByMeCount           int             `xml:"photosByMeCount"`
	PhotosWithMeCount         int             `xml:"photosWithMeCount"`
	CommentsEnabled           bool            `xml:"commentsEnabled" db:"comments_enabled"`
	Location                  common.Location `xml:"location"`
	LocationX                 int32           `xml:"-" db:"location_x"`
	LocationY                 int32           `xml:"-" db:"location_y"`
	FavouriteSlotCount        string          `xml:"favouriteslotcount"`
	FavouriteUserCount        string          `xml:"favouriteusercount"`
	QueuedLevelCount          string          `xml:"lolcatftwcount"`
	Pins                      string          `xml:"pins"`
	StaffChallengeGoldCount   string          `xml:"staffchallengegoldcount"`
	StaffChallengeSilverCount string          `xml:"staffchallengesilvercount"`
	StaffChallengeBronzeCount string          `xml:"staffchallengebronzecount"`
	ClientsConnected          struct {
		LittleBigPlanet     bool `xml:"lbp1"`
		LittleBigPlanet2    bool `xml:"lbp2"`
		LittleBigPlanet3PS4 bool `xml:"lbp3ps4"`
		LittleBigPlanet3PS3 bool `xml:"lbp3ps3"`
	} `xml:"clientsConnected"`

	BooHash string `xml:"boo2,omitempty" db:"boo_icon"`
	MehHash string `xml:"meh2,omitempty" db:"meh_icon"`
	YayHash string `xml:"yay2,omitempty" db:"yay_icon"`

	LevelVisibility   sql.NullString `xml:"-" db:"level_visibility"`
	ProfileVisibility sql.NullString `xml:"-" db:"profile_visibility"`
}
type PlanetUpdate struct {
	XMLName  xml.Name `xml:"user"`
	Planets  string   `xml:"planets,omitempty"`            // The resource of the planet for the current game
	CCPlanet string   `xml:"crossControlPlanet,omitempty"` // The planet for cross controller
}
type UpdateUser struct {
	XMLName   xml.Name
	Biography string                         `xml:"biography" db:"bio"`
	Icon      string                         `xml:"icon,omitempty"` // The resource of the users icon
	Location  *common.Location               `xml:"location,omitempty"`
	Slots     slot.List[slot.UpdateUserSlot] `xml:"slots,omitempty"`
	BooHash   string                         `xml:"boo2,omitempty"`
	MehHash   string                         `xml:"meh2,omitempty"`
	YayHash   string                         `xml:"yay2,omitempty"`
}
