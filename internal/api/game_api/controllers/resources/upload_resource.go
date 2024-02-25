package resources

import (
	db2 "HugeSpaceship/internal/hs_db"
	"HugeSpaceship/internal/model/auth"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/utils"
	"net/http"
)

func UploadResources() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := utils.GetContextValue[auth.Session](r.Context(), "session")
		hash := r.PathValue("hash")
		dbCtx := db.GetContext()

		exists, err := db2.ResourceExists(dbCtx, hash)
		if err != nil {
			ctx.Error(err)
			utils.HttpLog(w, http.StatusInternalServerError, "Failed to check if resource exists")
			return
		}

		// 2024 update: I have no idea what this comment is about, your guess is as good as mine.
		//oi im back but im not also uh i cant find the ps3! :(
		if exists {
			ctx.String(http.StatusConflict, "Resource already exists")
			return
		}

		err = db2.UploadResource(dbCtx, r.Body, r.ContentLength, hash, session.UserID)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatus(500)
			return
		}

		ctx.Status(200)
	}
}
