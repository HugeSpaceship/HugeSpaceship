package moderation

import "github.com/gin-gonic/gin"

// FilterHandler handles all text that can be moderated for the game
// Note: this is not used in LBPVita
func FilterHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		// This is a stub for now
		// TODO: implement text filtering
		context.DataFromReader(200, context.Request.ContentLength, "text/xml", context.Request.Body, nil)
	}
}
