package lbp_image

import (
	"HugeSpaceship/testdata"
	"bytes"
	"crypto/sha1"
	"io"
	"reflect"
	"testing"
)

// SHA1 Sum of the expected image
var ddsSum = []byte{0x44, 0x32, 0xc6, 0xa6, 0xe1, 0x18, 0x06, 0x2b, 0x6b, 0xd0, 0xa0, 0xb8, 0x1e, 0xa4, 0xd7, 0xb8, 0x0f, 0xc1, 0xfe, 0x72}
var ddsSize = int64(4736)

func TestDecompressImage(t *testing.T) {
	decompressedImage, err := DecompressImage(bytes.NewReader(testdata.TestDDSCompressed))
	if err != nil {
		t.Error(err)
	}
	shaSum := sha1.New()
	numCopied, err := io.Copy(shaSum, decompressedImage)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(shaSum.Sum(nil), ddsSum) {
		t.Error("sum does not match expected value")
	}
	if numCopied != ddsSize {
		t.Error("image is not the expected size")
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
		// TODO: Add testdata cases.
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
