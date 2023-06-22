package slots

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"github.com/gin-gonic/gin"
)

func GetSlotsByHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.XML(200, &lbp_xml.Slots{Total: 0, HintStart: 0})
	}
}
