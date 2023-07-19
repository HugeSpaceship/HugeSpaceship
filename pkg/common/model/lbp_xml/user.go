package lbp_xml

import (
	"encoding/xml"
)

type User struct {
	XMLName                   xml.Name `xml:"user"`
	Type                      string   `xml:"type,attr"`
	NpHandle                  NpHandle `xml:"npHandle"`
	Biography                 string   `xml:"biography"`
	Game                      string   `xml:"game"`
	Lbp1UsedSlots             string   `xml:"lbp1UsedSlots"`
	EntitledSlots             string   `xml:"entitledSlots"`
	FreeSlots                 string   `xml:"freeSlots"`
	CrossControlUsedSlots     string   `xml:"crossControlUsedSlots"`
	CrossControlEntitledSlots string   `xml:"crossControlEntitledSlots"`
	CrossControlFreeSlots     string   `xml:"crossControlFreeSlots"`
	Lbp2UsedSlots             string   `xml:"lbp2UsedSlots"`
	Lbp2EntitledSlots         string   `xml:"lbp2EntitledSlots"`
	Lbp2FreeSlots             string   `xml:"lbp2FreeSlots"`
	Lbp3UsedSlots             string   `xml:"lbp3UsedSlots"`
	Lbp3EntitledSlots         string   `xml:"lbp3EntitledSlots"`
	Lbp3FreeSlots             string   `xml:"lbp3FreeSlots"`
	Lists                     string   `xml:"lists"`
	ListsQuota                string   `xml:"lists_quota"`
	HeartCount                string   `xml:"heartCount"`
	ReviewCount               string   `xml:"reviewCount"`
	CommentCount              string   `xml:"commentCount"`
	PhotosByMeCount           string   `xml:"photosByMeCount"`
	PhotosWithMeCount         string   `xml:"photosWithMeCount"`
	CommentsEnabled           bool     `xml:"commentsEnabled"`
	Location                  Location `xml:"location"`
	FavouriteSlotCount        string   `xml:"favouriteslotcount"`
	FavouriteUserCount        string   `xml:"favouriteusercount"`
	QueuedLevelCount          string   `xml:"lolcatftwcount"`
	Pins                      string   `xml:"pins"`
	StaffChallengeGoldCount   string   `xml:"staffchallengegoldcount"`
	StaffChallengeSilverCount string   `xml:"staffchallengesilvercount"`
	StaffChallengeBronzeCount string   `xml:"staffchallengebronzecount"`
	ClientsConnected          struct {
		LittleBigPlanet     bool `xml:"lbp1"`
		LittleBigPlanet2    bool `xml:"lbp2"`
		LittleBigPlanet3PS4 bool `xml:"lbp3ps4"`
	} `xml:"clientsConnected"`
}
