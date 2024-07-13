package utils

type FileType string

const (
	Texture              = "TEX"
	PNG                  = "PNG"
	JPG                  = "JPG"
	Plan                 = "PLN"
	MotionRecording      = "REC"
	FishFingerScript     = "FSH"
	Sound                = "VOP"
	Level                = "LVL"
	Painting             = "PTG"
	AdventureCreate      = "ADC" // Adventure layout in earth/level list
	AdventureSharedData  = "ADS" // Shared data across an adventure (items, quests, start points, etc.)
	Quest                = "QST"
	StreamingChunk       = "CHK" // For dynamic thermometer levels
	CrossControllerLevel = "PRF"
	UNK                  = "UNK" // Placeholder for when we don't know
)
