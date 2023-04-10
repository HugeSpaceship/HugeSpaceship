package types

type DataType uint16
type SectionType uint8

const (
	Body   SectionType = 0x00
	Footer SectionType = 0x02
)

const (
	Empty     DataType = 0x00
	UInt32    DataType = 0x01
	UInt64    DataType = 0x02
	String    DataType = 0x04
	Timestamp DataType = 0x07
	Binary    DataType = 0x08
)

func TypeToString(dataType DataType) string {
	switch dataType {
	case Empty:
		return "Empty"
	case UInt32:
		return "UInt32"
	case UInt64:
		return "UInt64"
	case String:
		return "String"
	case Timestamp:
		return "Timestamp"
	case Binary:
		return "Binary"
	}
	return "Unknown"
}
