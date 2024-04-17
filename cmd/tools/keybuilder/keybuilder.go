// Key builder is a tool for creating PEM NpTicket keys

package main

import (
	"HugeSpaceship/pkg/npticket/signing"
	"HugeSpaceship/pkg/npticket/types"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log/slog"
	"math/big"
	"os"
)

var X = flag.String("x", "", "The X location on a ECDSA curve")
var Y = flag.String("y", "", "The Y location on a ECDSA curve")
var Curve = flag.String("curve", "secp192r1", "The ECDSA curve to use")
var Description = flag.String("desc", "A Key", "The description of the key, usually the platform it's from i.e. PSN, RPCN")
var Owner = flag.String("owner", "The Key Owner", "The entity who controls the signing of tickets using the public key")
var File = flag.String("file", "-", "The file to write the key to, or - for stdout")
var SignatoryID = flag.Uint("signatory", 123456, "The id of the npticket signatory")

func main() {
	flag.Parse()

	x, _ := new(big.Int).SetString(*X, 16)
	y, _ := new(big.Int).SetString(*Y, 16)

	if x == nil || x.Cmp(big.NewInt(0)) == 0 {
		slog.Error("X must not be a valid, non-zero hex number")
		os.Exit(1)
	}

	if y == nil || y.Cmp(big.NewInt(0)) == 0 {
		slog.Error("Y must not be a valid, non-zero hex number")
		os.Exit(1)
	}

	var key ecdsa.PublicKey
	if curve, ok := signing.CurveFromName(*Curve); ok {
		key.Curve = curve
	} else {
		slog.Error("Invalid Curve, must be one of: \nsecp224r1\nsecp256r1\nsecp384r1\nsecp521r1\nsecp192r1\nsecp224k1")
		os.Exit(1)
	}

	key.X = x
	key.Y = y

	info := types.KeyInfo{
		SignatoryID: uint32(*SignatoryID),
		Description: *Description,
		Owner:       *Owner,
	}
	keyBytes, err := signing.SaveKey(&key, info)
	if err != nil {
		slog.Error("failed to serialize key", slog.Any("error", err))
		os.Exit(1)
	}

	if *File == "-" {
		fmt.Println(string(keyBytes))
		return
	}

	err = os.WriteFile(*File, keyBytes, 0644)
	if err != nil {
		slog.Error("failed to write key", slog.Any("error", err))
	}
}
