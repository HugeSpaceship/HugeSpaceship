package auth

type LoginResult struct {
	AuthTicket string `xml:"authTicket"`
	LBPEnvVer  string `xml:"lbpEnvVer"`
}
