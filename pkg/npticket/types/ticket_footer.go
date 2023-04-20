package types

const (
	PSNSignatoryID  uint32 = 0x719F1D4A
	RPCNSignatoryID uint32 = 0x5250434E
)

type TicketFooter struct {
	Signatory uint32
	Signature []byte
}
