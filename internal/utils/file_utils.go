package utils

import (
	"bufio"
	"io"
)

func GetResourceType(r io.Reader) (io.Reader, FileType, error) {
	br := bufio.NewReader(r)
	magic, err := br.Peek(3)
	if err != nil {
		return br, Unknown, err
	}

	ft := Unknown
	err = ft.Scan(magic)

	return br, ft, err
}
