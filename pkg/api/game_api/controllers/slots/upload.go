package slots

import (
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func StartPublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slot := lbp_xml.Slot{}
		err := ctx.BindXML(&slot)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse xml body")
		}
		ctx.XML(200, lbp_xml.Slot{Resource: []string{slot.RootLevel}, Type: "user"})
	}
}

func PublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
