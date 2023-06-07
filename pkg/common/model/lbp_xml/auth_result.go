package lbp_xml

type AuthResult struct {
	AuthTicket string `xml:"authTicket"`
	LBPEnvVer  string `xml:"lbpEnvVer"`
}
