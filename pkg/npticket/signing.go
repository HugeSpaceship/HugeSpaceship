package npticket

import (
	"HugeSpaceship/pkg/npticket/types"
	"crypto/ecdsa"
	"encoding/asn1"
	"math/big"
	"sync"
)

var PSNKey *ecdsa.PublicKey
var RPCNKey *ecdsa.PublicKey

var initKeysOnce sync.Once

func initKeys() {
	var psnKey ecdsa.PublicKey
	psnKey.Curve = secp192r1()
	psnKey.X = big.NewInt(0)
	psnKey.Y = big.NewInt(0)
	psnKey.X.SetString("39c62d061d4ee35c5f3f7531de0af3cf918346526edac727", 16)
	psnKey.Y.SetString("a5d578b55113e612bf1878d4cc939d61a41318403b5bdf86", 16)
	PSNKey = &psnKey

	var rpcnKey ecdsa.PublicKey
	rpcnKey.Curve = S224()
	rpcnKey.X = big.NewInt(0)
	rpcnKey.Y = big.NewInt(0)
	rpcnKey.X.SetString("b07bc0f0addb97657e9f389039e8d2b9c97dc2a31d3042e7d0479b93", 16)
	rpcnKey.Y.SetString("d81c42b0abdf6c42191a31e31f93342f8f033bd529c2c57fdb5a0a7d", 16)
	RPCNKey = &rpcnKey
}

func VerifyTicket(ticket types.Ticket) bool {
	initKeysOnce.Do(initKeys)
	signature := ticket.Footer.Signature

	switch ticket.Footer.Signatory {
	case types.PSNSignatoryID:
		return ecdsa.VerifyASN1(PSNKey, ticket.BodyHash, signatureFix(signature))
	case types.RPCNSignatoryID:
		return ecdsa.VerifyASN1(RPCNKey, ticket.BodyHash, signature)
	}

	return false
}

type asn1Signature struct {
	R *big.Int
	S *big.Int
}

// signatureFix removes the extra zeroes from the end of a npticket signature
func signatureFix(signature []byte) []byte {
	for i := 0; i <= 2; i++ {
		remaining, err := asn1.Unmarshal(signature[:len(signature)-i], &asn1Signature{})
		if err == nil && len(remaining) == 0 {
			return signature[:len(signature)-i]
		}
	}
	return signature
}
