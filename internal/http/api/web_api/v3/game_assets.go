package v3

import (
	"bytes"
	resMan "github.com/HugeSpaceship/HugeSpaceship/internal/resources"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/image"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func GameAssetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func ImageAssetHandler(res *resMan.ResourceManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		resReader, err := res.GetResource(hash)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to access resource")
			return
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, resReader)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to read image")
			return
		}
		resReader.Close()

		decompressedImage, err := image.DecompressImage(buf)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to decompress image")
			return
		}
		outBuf := new(bytes.Buffer)
		err = image.IMGToPNG(decompressedImage, outBuf)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to convert image")
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(outBuf.Len()))
		_, err = io.Copy(w, outBuf)
		if err != nil {
			slog.Error("Failed to write response")
		}
	}
}
