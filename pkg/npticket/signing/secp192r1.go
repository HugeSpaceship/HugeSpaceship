package signing

import (
	"crypto/elliptic"
	"math/big"
)

var p192r1 *elliptic.CurveParams

func initSECP192R1() {
	p192r1 = &elliptic.CurveParams{Name: "P-192"}
	p192r1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFFFFFFFFFFFF", 16)
	p192r1.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFF99DEF836146BC9B1B4D22831", 16)
	p192r1.B, _ = new(big.Int).SetString("64210519E59C80E70FA7E9AB72243049FEB8DEECC146B9B1", 16)
	p192r1.Gx, _ = new(big.Int).SetString("188DA80EB03090F67CBF20EB43A18800F4FF0AFD82FF1012", 16)
	p192r1.Gy, _ = new(big.Int).SetString("07192B95FFC8DA78631011ED6B24CDD573F977A11E794811", 16)
	p192r1.BitSize = 192
}

func P192() elliptic.Curve {
	initonce.Do(initAll)
	return p192r1
}
