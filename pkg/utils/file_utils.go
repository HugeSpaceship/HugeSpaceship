package utils

import (
	"io"
	"os"
)

func ServeFile(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	return err
}
