// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/auth"
	"github.com/HugeSpaceship/HugeSpaceship/internal/resources"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"log/slog"
	"net/http"
)

func UploadResources(res *resources.ResourceManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := utils.GetContextValue[auth.Session](r.Context(), "session")
		hash := r.PathValue("hash")

		exists, err := res.HasResource(hash)
		if err != nil {
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to check if resource exists")
			return
		}

		//oi im back but im not also uh i cant find the ps3! :(
		// - SyngletOxygen 2023

		if exists {
			utils.HttpLog(w, http.StatusConflict, "Resource already exists")
			return
		}

		err = res.UploadResource(hash, r.Body, r.ContentLength, session.UserID)
		if err != nil {
			slog.Error("error saving resource", slog.Any("error", err))
			utils.HttpLog(w, http.StatusInternalServerError, "internal error in resource upload")
			return
		}
	}
}
