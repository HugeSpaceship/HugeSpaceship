package lbp_image

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	_ "github.com/hugespaceship/dds"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
)

var (
	// Magic numbers
	pngMagic  = []byte{0x89, 0x50, 0x4e, 0x47}
	jpegMagic = []byte{0x4A, 0x46, 0x49, 0x46}
	texMagic  = []byte("TEX ")
)

func decompressZlibData(reader io.Reader, len uint16) ([]byte, error) {
	zlibReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	inflatedData := make([]byte, len)
	n, err := io.ReadFull(zlibReader, inflatedData)

	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if n != int(len) {
		return nil, fmt.Errorf("%w got %d, expected %d", InvalidZlibLength, n, len)
	}
	return inflatedData, nil
}

// DecompressImage implements the necessary decompression to read an LBP texture.
// The code for this is shamelessly stolen from ProjectLighthouse, which in turn stole it from Toolbox
// Original code is here: https://github.com/ennuo/toolkit/blob/d996ee4134740db0ee94e2cbf1e4edbd1b5ec798/src/main/java/ennuo/craftworld/utilities/Compressor.java#L40
func DecompressImage(inReader io.Reader) (io.Reader, error) {
	reader := bufio.NewReader(inReader)

	magic, err := reader.Peek(9)
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(magic, pngMagic) || bytes.HasSuffix(magic, jpegMagic) {
		return reader, nil
	}

	if !bytes.HasPrefix(magic, texMagic) {
		return nil, fmt.Errorf("%w (%X)", InvalidMagicNumber, magic)
	}

	_, _ = reader.Discard(6) // skip the 4 byte magic number, and the two byte mystery number
	var chunks uint16
	err = binary.Read(reader, binary.BigEndian, &chunks)
	if err != nil {
		return nil, err
	}

	compressed := make([]uint16, chunks)
	decompressed := make([]uint16, chunks)

	for i := uint16(0); i < chunks; i++ {
		err := binary.Read(reader, binary.BigEndian, &compressed[i])
		if err != nil {
			return nil, err
		}
		err = binary.Read(reader, binary.BigEndian, &decompressed[i])
		if err != nil {
			return nil, err
		}
	}

	buf := new(bytes.Buffer)
	for i := uint16(0); i < chunks; i++ {
		deflatedData := make([]byte, compressed[i])
		_, err = io.ReadFull(reader, deflatedData)

		if compressed[i] == decompressed[i] {
			buf.Write(deflatedData)
			continue
		}

		inflatedData, err := decompressZlibData(bytes.NewReader(deflatedData), decompressed[i])
		if err != nil {
			return nil, err
		}

		buf.Write(inflatedData)
	}

	return bytes.NewReader(buf.Bytes()), nil
}

// IMGToPNG tries to convert any image to a PNG
func IMGToPNG(r io.Reader, w io.Writer) error {
	if r == nil {
		return errors.New("nil reader")
	}
	if w == nil {
		return errors.New("nil writer")
	}
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	return png.Encode(w, img)
}
