package resources

import "github.com/gin-gonic/gin"

func ShowModeratedHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//
		ctx.Data(200, "text/xml", []byte("<resources/>"))
	}
}
