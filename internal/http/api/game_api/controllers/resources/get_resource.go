package resources

import (
	"errors"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources/backends"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/validation"
	"io"
	"log/slog"
	"net/http"
)

func GetResourceHandler(res *resources.ResourceManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, hash := validation.IsHashValid(r.PathValue("hash"))
		if !ok {
			utils.HttpLog(w, http.StatusBadRequest, "Invalid resource hash")
			return
		}

		resReader, err := res.GetResource(r.Context(), hash)
		if err != nil {
			if errors.Is(err, backends.ResourceNotFound) {
				utils.HttpLog(w, http.StatusNotFound, "Resource not found")
				return
			}
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to get resource")
			slog.Error("Failed to open resource", slog.Any("err", err.Error()))
			return
		}
		defer func(resReader io.ReadCloser) {
			err := resReader.Close()
			if err != nil {
				panic(err)
			}
		}(resReader)

		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = io.Copy(w, resReader)
		if err != nil {
			slog.Error("Failed to write resource")
		}
	}
}
