package common

type GameType string

const (
	LBP1   GameType = "LBP1"   // LittleBigPlanet (2008)
	LBPPSP GameType = "LBPPSP" // LittleBigPlanet (2009)
	LBP2   GameType = "LBP2"   // LittleBigPlanet 2 (2011)
	LBPV   GameType = "LBPV"   // LittleBigPlanet PSVita (2012)
	LBP3   GameType = "LBP3"   // LittleBigPlanet 3 (2014)
)

func (g GameType) ToInt() int {
	switch g {
	case LBP1:
		return 0
	case LBP2:
		return 1
	case LBP3:
		return 2
	case LBPV:
		return 3
	case LBPPSP:
		return 4
	default:
		panic("invalid game type")
	}
}
