package types

import (
	"fmt"
)

type RoomSlot struct {
	Id   uint64
	Type string
}

func (s RoomSlot) IsNull() bool {
	return false
}

func (s RoomSlot) Index(i int) any {
	switch i {
	case 0:
		return s.Id
	case 1:
		return s.Type
	default:
		panic("invalid index")
	}
}

func (s *RoomSlot) ScanNull() error {
	return fmt.Errorf("cannot scan NULL into point3d")
}

func (s *RoomSlot) ScanIndex(i int) any {
	switch i {
	case 0:
		return &s.Id
	case 1:
		return &s.Type
	default:
		panic("invalid index")
	}
}
