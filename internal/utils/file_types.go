package utils

type FileType string

const (
	Texture              = "TEX" // DDS Texture
	PNG                  = "PNG" // PNG Texture
	JPG                  = "JPG" // JPEG Texture
	Plan                 = "PLN" // Object, from emitters, the popit, etc.
	MotionRecording      = "REC" // PSMove Motion Recording?
	FishFingerScript     = "FSH" // Executable script
	Sound                = "VOP" // Audio
	Level                = "LVL" // Level data
	Painting             = "PTG" // PSMove Painting?
	AdventureCreate      = "ADC" // Adventure layout in earth/level list
	AdventureSharedData  = "ADS" // Shared data across an adventure (items, quests, start points, etc.)
	Quest                = "QST" // LBP3 Quest
	StreamingChunk       = "CHK" // For dynamic thermometer levels
	CrossControllerLevel = "PRF" // Level for cross controller
	UNK                  = "UNK" // Placeholder for when we don't know
)
