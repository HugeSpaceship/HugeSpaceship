package utils

import (
	"fmt"
)

type FileType string

func (e *FileType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FileType(s)
	case string:
		*e = FileType(s)
	default:
		return fmt.Errorf("unsupported scan type for FileType: %T", src)
	}

	for _, fileType := range fileTypes { // Check that the FileType is valid
		if fileType == *e {
			return nil
		}
	}

	*e = Unknown
	return fmt.Errorf("invalid FileType: %T", src)
}

const (
	Texture              FileType = "TEX" // DDS Texture
	PNG                  FileType = "PNG" // PNG Texture
	JPG                  FileType = "JPG" // JPEG Texture
	Plan                 FileType = "PLN" // Object, from emitters, the popit, etc.
	MotionRecording      FileType = "REC" // PSMove Motion Recording?
	FishFingerScript     FileType = "FSH" // Executable script
	Sound                FileType = "VOP" // Audio
	Level                FileType = "LVL" // Level data
	Painting             FileType = "PTG" // PSMove Painting?
	AdventureCreate      FileType = "ADC" // Adventure layout in earth/level list
	AdventureSharedData  FileType = "ADS" // Shared data across an adventure (items, quests, start points, etc.)
	Quest                FileType = "QST" // LBP3 Quest
	StreamingChunk       FileType = "CHK" // For dynamic thermometer levels
	CrossControllerLevel FileType = "PRF" // Level for cross controller
	Unknown              FileType = "UNK" // Placeholder for when we don't know
)

// This is really bad, if someone can think of a better way of doing it then lmk
var fileTypes = []FileType{Texture, PNG, JPG, Plan, MotionRecording, FishFingerScript, Sound, Level, Painting, AdventureCreate, AdventureSharedData, Quest, StreamingChunk, CrossControllerLevel, Unknown}
