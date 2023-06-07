package users

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"github.com/gin-gonic/gin"
)

func UserGetHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.XML(200, &lbp_xml.ExampleUser)
	}
}
