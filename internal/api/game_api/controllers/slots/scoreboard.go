package slots

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

func UploadScoreHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	}
}
