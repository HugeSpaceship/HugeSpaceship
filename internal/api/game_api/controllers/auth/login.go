package auth

import (
	"HugeSpaceship/internal/api/game_api/utils"
	"HugeSpaceship/internal/hs_db/auth"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/npticket"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/netip"
)

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticketData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Err(err).Msg("failed to get request data")
		}
		parser := npticket.NewParser(ticketData)
		ticket, err := parser.Parse()

		log.Info().Str("userName", ticket.Username).Str("country", ticket.Country).Msg("User Connected")

		if !npticket.VerifyTicket(ticket) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		log.Debug().Msg("Verified NPTicket")
		game := r.URL.Query().Get("game")

		log.Debug().Msg("Getting Context")
		dbCtx := db.GetContext()
		defer db.CloseContext(dbCtx)
		log.Debug().Msg("Got Context, getting session")

		if r.Header.Get("X-exe-v") != "" {
			game = "lbp-psp"
		}

		token, err := auth.NewSession(dbCtx, ticket, netip.MustParseAddr(r.RemoteAddr), game, c.Query("titleID"))
		if err != nil {
			c.AbortWithStatus(403)
			return
		}

		if c.GetHeader("X-exe-v") == "" { // if we're not on PSP
			c.Render(200, utils.LBPXML{Data: &lbp_xml.LoginResult{
				AuthTicket: "MM_AUTH=" + token,
				LBPEnvVer:  "HugeSpaceship",
			}})
		} else {
			c.Render(200, utils.LBPXML{Data: &lbp_xml.PSPLoginResult{
				AuthTicket: "MM_AUTH=" + token,
			}})
		}
	}
}
