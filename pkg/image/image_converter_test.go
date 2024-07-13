package image

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/HugeSpaceship/HugeSpaceship/testdata"
	"io"
	"os"
	"reflect"
	"testing"
)

// SHA1 Sum of the expected image
var ddsSum = []byte{0x44, 0x32, 0xc6, 0xa6, 0xe1, 0x18, 0x06, 0x2b, 0x6b, 0xd0, 0xa0, 0xb8, 0x1e, 0xa4, 0xd7, 0xb8, 0x0f, 0xc1, 0xfe, 0x72}
var ddsSize = int64(4736)

func mustBytes(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return b
}

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
		name     string
		args     args
		wantHash string
		wantErr  bool
	}{
		// TODO: Add testdata cases.
		{
			name: "valid image",
			args: args{
				r: bytes.NewBuffer(mustBytes(os.ReadFile("../../testdata/image/test.dds"))),
			},
			wantHash: "f058a07ad5c855c41261b301aa3235a68db3be1481c419f9a1b30ed441053d24",
			wantErr:  false,
		},
		{
			name: "nil reader",
			args: args{
				r: nil,
			},
			wantErr: true,
		},
		{
			name: "empty image",
			args: args{
				r: new(bytes.Buffer),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := sha256.New()
			err := IMGToPNG(tt.args.r, hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("IMGToPNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gotHash := hex.EncodeToString(hash.Sum(nil))
				if gotHash != tt.wantHash {
					t.Errorf("IMGToPNG() gotHash = %0x, want %0x", gotHash, tt.wantHash)
				}
			}
		})
	}
}
