package signing

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509/pkix"
	_ "embed"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/HugeSpaceship/HugeSpaceship/pkg/npticket/types"
	"log/slog"
	"math/big"
	"strconv"
	"sync"
)

var initKeysOnce sync.Once

const NpTicketKeyType = "NPTICKET PUBLIC KEY"
const SignatoryIDHeader = "Signatory-ID"
const DescriptionHeader = "Description"
const OwnerHeader = "Owner"

var ErrorInvalidCurve = errors.New("invalid curve")

type asn1PublicKey struct {
	Algo      pkix.AlgorithmIdentifier
	BitString asn1.BitString
}

func SaveKey(key *ecdsa.PublicKey, info types.KeyInfo) ([]byte, error) {
	data := elliptic.Marshal(key.Curve, key.X, key.Y)

	curveOID, ok := oidFromNamedCurve(key.Curve)
	if !ok {
		return nil, ErrorInvalidCurve
	}
	asnKey := asn1PublicKey{
		Algo: pkix.AlgorithmIdentifier{
			Algorithm: curveOID,
		},
		BitString: asn1.BitString{
			Bytes:     data,
			BitLength: 8 * len(data),
		},
	}

	asnData, err := asn1.Marshal(asnKey)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = pem.Encode(buf, &pem.Block{
		Type: NpTicketKeyType,
		Headers: map[string]string{
			SignatoryIDHeader: strconv.Itoa(int(info.SignatoryID)),
			OwnerHeader:       info.Owner,
			DescriptionHeader: info.Description,
		},
		Bytes: asnData,
	})

	return buf.Bytes(), err
}

var ErrorPemKeyNotFound = errors.New("pem key not found")
var ErrorInvalidKeyType = errors.New("invalid key type")
var ErrorMissingSignatoryID = errors.New("missing signatory ID")
var ErrorInvalidSignatoryID = errors.New("invalid signatory ID")

func LoadKey(data []byte) (*types.NpTicketPublicKey, []byte, error) {
	block, rest := pem.Decode(data)
	if block == nil {
		return nil, nil, ErrorPemKeyNotFound
	}

	if block.Type != NpTicketKeyType {
		return nil, nil, ErrorInvalidKeyType
	}

	if _, exists := block.Headers[SignatoryIDHeader]; !exists {
		return nil, nil, ErrorMissingSignatoryID
	}

	signatoryID, err := strconv.ParseUint(block.Headers[SignatoryIDHeader], 10, 32)
	if err != nil {
		return nil, nil, ErrorInvalidSignatoryID
	}
	info := types.KeyInfo{
		SignatoryID: uint32(signatoryID),
	}
	if owner, exists := block.Headers[OwnerHeader]; exists {
		info.Owner = owner
	}
	if desc, exists := block.Headers[DescriptionHeader]; exists {
		info.Description = desc
	}

	asnKey := new(asn1PublicKey)
	_, err = asn1.Unmarshal(block.Bytes, asnKey)
	if err != nil {
		return nil, nil, err
	}
	curve, ok := curveFromOid(asnKey.Algo.Algorithm)
	if !ok {
		return nil, nil, ErrorInvalidCurve
	}

	x, y := elliptic.Unmarshal(curve, asnKey.BitString.Bytes)

	if x == nil {
		return nil, nil, ErrorInvalidCurve
	}

	key := &types.NpTicketPublicKey{
		Info: &info,
		Key: &ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
	}
	return key, rest, nil
}

//go:embed keychain.pem
var builtinKeychain []byte

var keychain map[uint32]*types.NpTicketPublicKey

func initKeys() {
	keychain = make(map[uint32]*types.NpTicketPublicKey)

	key, rest, err := LoadKey(builtinKeychain)
	if err != nil {
		panic(err)
	}
	keychain[key.Info.SignatoryID] = key

	for len(rest) > 0 {
		key, rest, err = LoadKey(rest)
		if err != nil {
			slog.Debug("Error loading key, however this is the builtin keys, so not to worry",
				slog.Any("error", err),
			)
			break
		}
		keychain[key.Info.SignatoryID] = key
	}
}

func VerifyTicket(ticket types.Ticket) bool {
	initKeysOnce.Do(initKeys)
	signature := ticket.Footer.Signature

	key, ok := keychain[ticket.Footer.Signatory]
	if !ok {
		return false
	}

	fmt.Printf("Ticket from %s signed by %s (%s)\n", ticket.Username, key.Info.Owner, key.Info.Description)

	if ecdsa.VerifyASN1(key.Key, ticket.BodyHash, signatureFix(signature)) {
		return true
	}

	return ecdsa.VerifyASN1(key.Key, ticket.BodyHash, signature)
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
