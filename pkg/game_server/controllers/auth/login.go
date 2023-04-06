package auth

import (
	"HugeSpaceship/pkg/model/auth"
	"HugeSpaceship/pkg/npticket"
	"github.com/gin-gonic/gin"
	"log"
)

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ticket, err := c.GetRawData()
		if err != nil {
			log.Println(err.Error())
		}
		npticket.Parse(ticket)
		c.XML(200, auth.LoginResult{
			AuthTicket: "testToken", // TODO: get real token
			LBPEnvVer:  "HugeSpaceship",
		})
	}
}
