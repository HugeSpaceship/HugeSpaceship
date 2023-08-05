package lbp_xml

import (
	"HugeSpaceship/pkg/common/model/common"
	"HugeSpaceship/pkg/common/model/lbp_xml/npdata"
	"encoding/xml"
	"github.com/google/uuid"
)

type User struct {
	XMLName                   xml.Name        `xml:"user"`
	ID                        uuid.UUID       `xml:"-" db:"id"`
	PSN_UID                   string          `xml:"-" db:"psn_uid"`
	RPCN_UID                  string          `xml:"-" db:"rpcn_uid"`
	Type                      string          `xml:"type,attr"`
	NpHandle                  npdata.NpHandle `xml:"npHandle"`
	Username                  string          `xml:"-" db:"username"`
	AvatarHash                string          `xml:"-" db:"avatar_hash"`
	Biography                 string          `xml:"biography" db:"bio"`
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
	LocationX                 int32           `xml:"-" db:"locationx"`
	LocationY                 int32           `xml:"-" db:"locationy"`
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
	} `xml:"clientsConnected"`
}
