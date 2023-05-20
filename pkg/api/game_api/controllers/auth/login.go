package auth

import (
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/npticket"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

		c.XML(200, auth.LoginResult{
			AuthTicket: "testToken", // TODO: get real token
			LBPEnvVer:  "HugeSpaceship",
		})
	}
}
