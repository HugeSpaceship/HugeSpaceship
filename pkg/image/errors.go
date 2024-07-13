package image

import (
	"errors"
)

var InvalidMagicNumber = errors.New("invalid magic number")

var InvalidZlibLength = errors.New("invalid length in zlib stream")
