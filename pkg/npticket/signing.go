package npticket

import (
	"crypto/ecdsa"
	"crypto/sha1"
	"math/big"
)

func (parser TicketParser) VerifyTicket() bool {
	var key ecdsa.PublicKey
	key.Curve = secp192r1()
	key.X = big.NewInt(0)
	key.Y = big.NewInt(0)
	key.X.SetString("39c62d061d4ee35c5f3f7531de0af3cf918346526edac727", 16)
	key.Y.SetString("a5d578b55113e612bf1878d4cc939d61a41318403b5bdf86", 16)

	hash := sha1.Sum(parser.TicketBody)
	return ecdsa.VerifyASN1(&key, hash[:], parser.ticket.Footer.Signature)
}
