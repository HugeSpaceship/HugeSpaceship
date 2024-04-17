package signing

import (
	"crypto/elliptic"
	"encoding/asn1"
)

var (
	oidNamedCurveP224   = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	oidNamedCurveP256   = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	oidNamedCurveP384   = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	oidNamedCurveP521   = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
	oidNamedCurveP192   = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 1}
	oidNamedCurveP224K1 = asn1.ObjectIdentifier{1, 3, 132, 0, 32}
)

func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) {
	switch curve {
	case elliptic.P224():
		return oidNamedCurveP224, true
	case elliptic.P256():
		return oidNamedCurveP256, true
	case elliptic.P384():
		return oidNamedCurveP384, true
	case elliptic.P521():
		return oidNamedCurveP521, true
	case P192():
		return oidNamedCurveP192, true
	case S224():
		return oidNamedCurveP224K1, true
	}

	return nil, false
}

var curveMap = map[string]elliptic.Curve{
	"1.3.132.0.33":        elliptic.P224(),
	"1.2.840.10045.3.1.7": elliptic.P256(),
	"1.3.132.0.34":        elliptic.P384(),
	"1.3.132.0.35":        elliptic.P521(),
	"1.2.840.10045.3.1.1": P192(),
	"1.3.132.0.32":        S224(),
}
var curveNameMap = map[string]elliptic.Curve{
	"secp224r1": elliptic.P224(),
	"secp256r1": elliptic.P256(),
	"secp384r1": elliptic.P384(),
	"secp521r1": elliptic.P521(),
	"secp192r1": P192(),
	"secp224k1": S224(),
}

func curveFromOid(identifier asn1.ObjectIdentifier) (elliptic.Curve, bool) {
	curve, ok := curveMap[identifier.String()]
	return curve, ok
}

func CurveFromName(name string) (elliptic.Curve, bool) {
	curve, ok := curveNameMap[name]
	return curve, ok
}
