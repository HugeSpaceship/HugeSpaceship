package image_utils

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	_ "github.com/lukegb/dds"
	"github.com/rs/zerolog/log"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
)

// DecompressImage implements the necessary decompression to read an LBP texture.
// The code for this is shamelessly stolen from ProjectLighthouse, which in turn stole it from Toolbox
// Original code is here: https://github.com/ennuo/toolkit/blob/d996ee4134740db0ee94e2cbf1e4edbd1b5ec798/src/main/java/ennuo/craftworld/utilities/Compressor.java#L40
func DecompressImage(closer io.ReadCloser) io.Reader {
	reader := bufio.NewReader(closer)
	_, _ = reader.Discard(3)

	if readMethod, _, _ := reader.ReadRune(); readMethod != ' ' {
		log.Error().Msg("Invalid image data")
		return nil
	}

	_, _ = reader.Discard(2)
	var chunks uint16
	err := binary.Read(reader, binary.BigEndian, &chunks)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode image")
		return nil
	}

	compressed := make([]uint16, chunks)
	decompressed := make([]uint16, chunks)

	for i := uint16(0); i < chunks; i++ {
		err := binary.Read(reader, binary.BigEndian, &compressed[i])
		if err != nil {
			log.Error().Err(err).Msg("Failed to read chunk header")
		}
		err = binary.Read(reader, binary.BigEndian, &decompressed[i])
		if err != nil {
			log.Error().Err(err).Msg("Failed to read chunk header")
		}
	}

	writer := new(bytes.Buffer)
	for i := uint16(0); i < chunks; i++ {
		deflatedData := make([]byte, compressed[i])
		_, err = reader.Read(deflatedData)

		if compressed[i] == decompressed[i] {
			writer.Write(deflatedData)
			continue
		}

		zlibReader, err := zlib.NewReader(bytes.NewReader(deflatedData))
		if err != nil {
			log.Error().Err(err).Msg("Failed to create zlib reader")
		}
		defer zlibReader.Close()
		inflatedData := make([]byte, decompressed[i])
		_, err = zlibReader.Read(inflatedData)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read compressed chunk")
		}

		writer.Write(inflatedData)
	}
	return bytes.NewReader(writer.Bytes())
}

// IMGToPNG tries to convert any image to a PNG
func IMGToPNG(r io.Reader, w io.Writer) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	return png.Encode(w, img)
}
