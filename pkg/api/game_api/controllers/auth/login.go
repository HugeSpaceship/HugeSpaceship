package auth

import (
	"HugeSpaceship/pkg/common/db/auth"
	"HugeSpaceship/pkg/common/model/lbp_xml"
	"HugeSpaceship/pkg/npticket"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/netip"
)

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ticketData, err := c.GetRawData()
		if err != nil {
			log.Err(err).Msg("failed to get request data")
		}
		parser := npticket.NewParser(ticketData)
		ticket, err := parser.Parse()

		log.Info().Str("userName", ticket.Username).Str("country", ticket.Country).Msg("User Connected")

		if !npticket.VerifyTicket(ticket) {
			c.Status(403)
			return
		}

		c.XML(200, lbp_xml.AuthResult{
			AuthTicket: "MM_AUTH=" + auth.NewSession(ticket, netip.MustParseAddr(c.ClientIP())), // TODO: get real token
			LBPEnvVer:  "HugeSpaceship",
		})
	}
}
