package lbp_xml

import "encoding/xml"

type LoginResult struct {
	XMLName         xml.Name `xml:"loginResult"`
	AuthTicket      string   `xml:"authTicket"`
	LBPEnvVer       string   `xml:"lbpEnvVer"`
	TitleStorageURL string   `xml:"titleStorageURL"`
}

// PSPLoginResult is the special login result for the PSP,
// as it uses a different response from every other LBP game
type PSPLoginResult struct {
	XMLName    xml.Name `xml:"authTicket"`
	AuthTicket string   `xml:",innerxml"`
}
