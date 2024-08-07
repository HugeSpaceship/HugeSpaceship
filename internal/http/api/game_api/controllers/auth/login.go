package auth

import (
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/auth"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml"
	utils2 "github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/npticket"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/npticket/signing"
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

		if !signing.VerifyTicket(ticket) {
			w.WriteHeader(http.StatusForbidden)
			slog.Debug("Rejecting npticket, as it is invalid")
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
			slog.Error("Failed to create session", "error", err)
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
