package types

const (
	PSNSignatoryID uint32 = 0x719F1D4A
)

type TicketFooter struct {
	Signatory uint32
	Signature []byte
}
