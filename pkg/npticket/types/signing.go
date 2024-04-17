package types

import "crypto/ecdsa"

type KeyInfo struct {
	SignatoryID uint32
	Owner       string
	Description string
}

type NpTicketPublicKey struct {
	Info *KeyInfo
	Key  *ecdsa.PublicKey
}
