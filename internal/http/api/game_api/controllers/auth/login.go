package auth

import (
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db"
	"github.com/HugeSpaceship/HugeSpaceship/internal/db/auth"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml"
	"github.com/HugeSpaceship/HugeSpaceship/internal/utils"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/npticket"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/npticket/signing"
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
			slog.Error("Failed to read body", "error", err)
		}
		parser := npticket.NewParser(ticketData)
		ticket, err := parser.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("Failed to parse ticket", "error", err)
		}

		slog.Debug("User Connected", "userName", ticket.Username, "country", ticket.Country)

		if !signing.VerifyTicket(ticket) {
			w.WriteHeader(http.StatusForbidden)
			slog.Info("Rejecting NpTicket with invalid signature")
			return
		}

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
			utils.HttpLog(w, http.StatusForbidden, "Failed to create session")
			slog.Error("Failed to create session", "error", err)
			return
		}

		w.Header().Set("Content-Type", "text/xml")

		if r.Header.Get("X-exe-v") == "" { // if we're not on PSP
			fmt.Println("Sending PS3 Login Result")
			err = utils.XMLMarshal(w, &lbp_xml.LoginResult{
				AuthTicket:      "MM_AUTH=" + token,
				LBPEnvVer:       "HugeSpaceship",
				TitleStorageURL: "http://dev.hugespaceship.io/api/LBP_XML",
			})
			if err != nil {
				slog.Error("Failed to marshal XML", slog.Any("error", err))
			}
		} else {
			err = utils.XMLMarshal(w, &lbp_xml.PSPLoginResult{
				AuthTicket: "MM_AUTH=" + token,
			})
			if err != nil {
				slog.Error("Failed to marshal XML", slog.Any("error", err))
			}
		}
	}
}
