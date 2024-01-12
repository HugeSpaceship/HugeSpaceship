package image_utils

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestDecompressImage(t *testing.T) {
	f, err := os.Open("../../test/b6d6869e139023340dc1f5225c1112e83acbbdca")
	if err != nil {
		t.Fatal(err)
	}
	decompressedImage := DecompressImage(f)
	f2, err := os.OpenFile("../../test/test.png", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(f2, decompressedImage)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIMGToPNG(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := IMGToPNG(tt.args.r, w)
			if (err != nil) != tt.wantErr {
				t.Errorf("IMGToPNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("IMGToPNG() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
