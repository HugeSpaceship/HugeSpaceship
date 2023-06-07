package lbp_xml

import "encoding/xml"

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
	Location                  struct {
		X int `xml:"x"`
		Y int `xml:"y"`
	}
	FavouriteSlotCount        string `xml:"favouriteslotcount"`
	FavouriteUserCount        string `xml:"favouriteusercount"`
	QueuedLevelCount          string `xml:"lolcatftwcount"`
	Pins                      string `xml:"pins"`
	StaffChallengeGoldCount   string `xml:"staffchallengegoldcount"`
	StaffChallengeSilverCount string `xml:"staffchallengesilvercount"`
	StaffChallengeBronzeCount string `xml:"staffchallengebronzecount"`
	ClientsConnected          struct {
		LittleBigPlanet     bool `xml:"lbp1"`
		LittleBigPlanet2    bool `xml:"lbp2"`
		LittleBigPlanet3PS4 bool `xml:"lbp3ps4"`
	} `xml:"clientsConnected"`
}

var ExampleUser = User{
	Type: "user",
	NpHandle: NpHandle{
		Username: "Zaprit282",
		IconHash: "g29838",
	},
	Game:                      "2",
	Lbp1UsedSlots:             "0",
	EntitledSlots:             "20",
	FreeSlots:                 "20",
	CrossControlUsedSlots:     "0",
	CrossControlEntitledSlots: "20",
	CrossControlFreeSlots:     "20",
	Lbp2UsedSlots:             "0",
	Lbp2EntitledSlots:         "20",
	Lbp2FreeSlots:             "20",
	Lbp3UsedSlots:             "0",
	Lbp3EntitledSlots:         "20",
	Lbp3FreeSlots:             "20",
	Lists:                     "0",
	ListsQuota:                "0",
	HeartCount:                "69",
	ReviewCount:               "2",
	CommentCount:              "0",
	PhotosByMeCount:           "1",
	PhotosWithMeCount:         "4",
	CommentsEnabled:           true,
	Location: struct {
		X int `xml:"x"`
		Y int `xml:"y"`
	}{
		X: 21068,
		Y: 30140,
	},
	FavouriteSlotCount:        "2",
	FavouriteUserCount:        "0",
	QueuedLevelCount:          "0",
	Pins:                      "1600196991,7235671,2051925987",
	StaffChallengeGoldCount:   "0",
	StaffChallengeSilverCount: "0",
	StaffChallengeBronzeCount: "0",
	ClientsConnected: struct {
		LittleBigPlanet     bool `xml:"lbp1"`
		LittleBigPlanet2    bool `xml:"lbp2"`
		LittleBigPlanet3PS4 bool `xml:"lbp3ps4"`
	}{
		LittleBigPlanet:     true,
		LittleBigPlanet2:    true,
		LittleBigPlanet3PS4: true,
	},
}
