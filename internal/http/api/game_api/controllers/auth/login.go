package auth

import (
	"HugeSpaceship/internal/hs_db/auth"
	"HugeSpaceship/internal/model/lbp_xml"
	"HugeSpaceship/pkg/db"
	"HugeSpaceship/pkg/npticket"
	utils2 "HugeSpaceship/pkg/utils"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"log/slog"
	"net/http"
	"net/netip"
	"strings"
)

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticketData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Err(err).Msg("failed to get request data")
		}
		parser := npticket.NewParser(ticketData)
		ticket, err := parser.Parse()

		log.Debug().Str("userName", ticket.Username).Str("country", ticket.Country).Msg("User Connected")

		if !npticket.VerifyTicket(ticket) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		log.Debug().Msg("Verified NPTicket")
		game := r.URL.Query().Get("game")

		conn, err := db.GetRequestConnection(r)
		if err != nil {
			panic(err)
		}

		if r.Header.Get("X-exe-v") != "" {
			game = "lbp-psp"
		}

		token, err := auth.NewSession(conn, ticket, netip.MustParseAddr(strings.Split(r.RemoteAddr, ":")[0]), game, r.URL.Query().Get("titleID"))
		if err != nil {
			utils2.HttpLog(w, http.StatusForbidden, "Failed to create session")
			return
		}

		w.Header().Set("Content-Type", "text/xml")

		if r.Header.Get("X-exe-v") == "" { // if we're not on PSP
			fmt.Println("Sending PS3 Login Result")
			err = utils2.XMLMarshal(w, &lbp_xml.LoginResult{
				AuthTicket:      "MM_AUTH=" + token,
				LBPEnvVer:       "HugeSpaceship",
				TitleStorageURL: "http://dev.hugespaceship.io/api/LBP_XML",
			})
			if err != nil {
				slog.Error("Failed to marshal XML", slog.Any("error", err))
			}
		} else {
			err = utils2.XMLMarshal(w, &lbp_xml.PSPLoginResult{
				AuthTicket: "MM_AUTH=" + token,
			})
			if err != nil {
				slog.Error("Failed to marshal XML", slog.Any("error", err))
			}
		}
	}
}
