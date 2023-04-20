package npticket

import (
	"crypto/elliptic"
	"math/big"
)

// A BitCurve represents a Koblitz Curve with a=0.
// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw.html
type BitCurve struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	B       *big.Int // the constant of the BitCurve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
}

func (BitCurve *BitCurve) Params() *elliptic.CurveParams {
	return &elliptic.CurveParams{
		P:       BitCurve.P,
		N:       BitCurve.N,
		B:       BitCurve.B,
		Gx:      BitCurve.Gx,
		Gy:      BitCurve.Gy,
		BitSize: BitCurve.BitSize,
	}
}

// IsOnCurve returns true if the given (x,y) lies on the BitCurve.
func (BitCurve *BitCurve) IsOnCurve(x, y *big.Int) bool {
	// y² = x³ + b
	y2 := new(big.Int).Mul(y, y) //y²
	y2.Mod(y2, BitCurve.P)       //y²%P

	x3 := new(big.Int).Mul(x, x) //x²
	x3.Mul(x3, x)                //x³

	x3.Add(x3, BitCurve.B) //x³+B
	x3.Mod(x3, BitCurve.P) //(x³+B)%P

	return x3.Cmp(y2) == 0
}

// TODO: double check if the function is okay
// affineFromJacobian reverses the Jacobian transform. See the comment at the
// top of the file.
func (BitCurve *BitCurve) affineFromJacobian(x, y, z *big.Int) (xOut, yOut *big.Int) {
	zinv := new(big.Int).ModInverse(z, BitCurve.P)
	zinvsq := new(big.Int).Mul(zinv, zinv)

	xOut = new(big.Int).Mul(x, zinvsq)
	xOut.Mod(xOut, BitCurve.P)
	zinvsq.Mul(zinvsq, zinv)
	yOut = new(big.Int).Mul(y, zinvsq)
	yOut.Mod(yOut, BitCurve.P)
	return
}

// Add returns the sum of (x1,y1) and (x2,y2)
func (BitCurve *BitCurve) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	z := new(big.Int).SetInt64(1)
	return BitCurve.affineFromJacobian(BitCurve.addJacobian(x1, y1, z, x2, y2, z))
}

// addJacobian takes two points in Jacobian coordinates, (x1, y1, z1) and
// (x2, y2, z2) and returns their sum, also in Jacobian form.
func (BitCurve *BitCurve) addJacobian(x1, y1, z1, x2, y2, z2 *big.Int) (*big.Int, *big.Int, *big.Int) {
	// See http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1 := new(big.Int).Mul(z1, z1)
	z1z1.Mod(z1z1, BitCurve.P)
	z2z2 := new(big.Int).Mul(z2, z2)
	z2z2.Mod(z2z2, BitCurve.P)

	u1 := new(big.Int).Mul(x1, z2z2)
	u1.Mod(u1, BitCurve.P)
	u2 := new(big.Int).Mul(x2, z1z1)
	u2.Mod(u2, BitCurve.P)
	h := new(big.Int).Sub(u2, u1)
	if h.Sign() == -1 {
		h.Add(h, BitCurve.P)
	}
	i := new(big.Int).Lsh(h, 1)
	i.Mul(i, i)
	j := new(big.Int).Mul(h, i)

	s1 := new(big.Int).Mul(y1, z2)
	s1.Mul(s1, z2z2)
	s1.Mod(s1, BitCurve.P)
	s2 := new(big.Int).Mul(y2, z1)
	s2.Mul(s2, z1z1)
	s2.Mod(s2, BitCurve.P)
	r := new(big.Int).Sub(s2, s1)
	if r.Sign() == -1 {
		r.Add(r, BitCurve.P)
	}
	r.Lsh(r, 1)
	v := new(big.Int).Mul(u1, i)

	x3 := new(big.Int).Set(r)
	x3.Mul(x3, x3)
	x3.Sub(x3, j)
	x3.Sub(x3, v)
	x3.Sub(x3, v)
	x3.Mod(x3, BitCurve.P)

	y3 := new(big.Int).Set(r)
	v.Sub(v, x3)
	y3.Mul(y3, v)
	s1.Mul(s1, j)
	s1.Lsh(s1, 1)
	y3.Sub(y3, s1)
	y3.Mod(y3, BitCurve.P)

	z3 := new(big.Int).Add(z1, z2)
	z3.Mul(z3, z3)
	z3.Sub(z3, z1z1)
	if z3.Sign() == -1 {
		z3.Add(z3, BitCurve.P)
	}
	z3.Sub(z3, z2z2)
	if z3.Sign() == -1 {
		z3.Add(z3, BitCurve.P)
	}
	z3.Mul(z3, h)
	z3.Mod(z3, BitCurve.P)

	return x3, y3, z3
}

// Double returns 2*(x,y)
func (BitCurve *BitCurve) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	z1 := new(big.Int).SetInt64(1)
	return BitCurve.affineFromJacobian(BitCurve.doubleJacobian(x1, y1, z1))
}

// doubleJacobian takes a point in Jacobian coordinates, (x, y, z), and
// returns its double, also in Jacobian form.
func (BitCurve *BitCurve) doubleJacobian(x, y, z *big.Int) (*big.Int, *big.Int, *big.Int) {
	// See http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l

	a := new(big.Int).Mul(x, x) //X1²
	b := new(big.Int).Mul(y, y) //Y1²
	c := new(big.Int).Mul(b, b) //B²

	d := new(big.Int).Add(x, b) //X1+B
	d.Mul(d, d)                 //(X1+B)²
	d.Sub(d, a)                 //(X1+B)²-A
	d.Sub(d, c)                 //(X1+B)²-A-C
	d.Mul(d, big.NewInt(2))     //2*((X1+B)²-A-C)

	e := new(big.Int).Mul(big.NewInt(3), a) //3*A
	f := new(big.Int).Mul(e, e)             //E²

	x3 := new(big.Int).Mul(big.NewInt(2), d) //2*D
	x3.Sub(f, x3)                            //F-2*D
	x3.Mod(x3, BitCurve.P)

	y3 := new(big.Int).Sub(d, x3)                  //D-X3
	y3.Mul(e, y3)                                  //E*(D-X3)
	y3.Sub(y3, new(big.Int).Mul(big.NewInt(8), c)) //E*(D-X3)-8*C
	y3.Mod(y3, BitCurve.P)

	z3 := new(big.Int).Mul(y, z) //Y1*Z1
	z3.Mul(big.NewInt(2), z3)    //3*Y1*Z1
	z3.Mod(z3, BitCurve.P)

	return x3, y3, z3
}

// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (BitCurve *BitCurve) ScalarMult(Bx, By *big.Int, k []byte) (*big.Int, *big.Int) {
	// We have a slight problem in that the identity of the group (the
	// point at infinity) cannot be represented in (x, y) form on a finite
	// machine. Thus the standard add/double algorithm has to be tweaked
	// slightly: our initial state is not the identity, but x, and we
	// ignore the first true bit in |k|.  If we don't find any true bits in
	// |k|, then we return nil, nil, because we cannot return the identity
	// element.

	Bz := new(big.Int).SetInt64(1)
	x := Bx
	y := By
	z := Bz

	seenFirstTrue := false
	for _, thingyByte := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			if seenFirstTrue {
				x, y, z = BitCurve.doubleJacobian(x, y, z)
			}
			if thingyByte&0x80 == 0x80 {
				if !seenFirstTrue {
					seenFirstTrue = true
				} else {
					x, y, z = BitCurve.addJacobian(Bx, By, Bz, x, y, z)
				}
			}
			thingyByte <<= 1
		}
	}

	if !seenFirstTrue {
		return nil, nil
	}

	return BitCurve.affineFromJacobian(x, y, z)
}

// ScalarBaseMult returns k*G, where G is the base point of the group and k is
// an integer in big-endian form.
func (BitCurve *BitCurve) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	return BitCurve.ScalarMult(BitCurve.Gx, BitCurve.Gy, k)
}

// Marshal converts a point into the form specified in section 4.3.6 of ANSI
// X9.62.
func (BitCurve *BitCurve) Marshal(x, y *big.Int) []byte {
	byteLen := (BitCurve.BitSize + 7) >> 3
	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // uncompressed point
	xBytes := x.Bytes()
	copy(ret[1+byteLen-len(xBytes):], xBytes)
	yBytes := y.Bytes()
	copy(ret[1+2*byteLen-len(yBytes):], yBytes)
	return ret
}

// Unmarshal converts a point, serialised by Marshal, into an x, y pair. On
// error, x = nil.
func (BitCurve *BitCurve) Unmarshal(data []byte) (x, y *big.Int) {
	byteLen := (BitCurve.BitSize + 7) >> 3
	if len(data) != 1+2*byteLen {
		return
	}
	if data[0] != 4 { // uncompressed form
		return
	}
	x = new(big.Int).SetBytes(data[1 : 1+byteLen])
	y = new(big.Int).SetBytes(data[1+byteLen:])
	return
}

//curve parameters taken from:
//http://www.secg.org/collateral/sec2_final.pdf

var secp160k1 *BitCurve
var secp192k1 *BitCurve
var secp224k1 *BitCurve
var secp256k1 *BitCurve

func initAll() {
	initS224()
	initSECP192R1()
}

func initS224() {
	// See SEC 2 section 2.6.1
	secp224k1 = new(BitCurve)
	secp224k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFE56D", 16)
	secp224k1.N, _ = new(big.Int).SetString("010000000000000000000000000001DCE8D2EC6184CAF0A971769FB1F7", 16)
	secp224k1.B, _ = new(big.Int).SetString("00000000000000000000000000000000000000000000000000000005", 16)
	secp224k1.Gx, _ = new(big.Int).SetString("A1455B334DF099DF30FC28A169A467E9E47075A90F7E650EB6B7A45C", 16)
	secp224k1.Gy, _ = new(big.Int).SetString("7E089FED7FBA344282CAFBD6F7E319F7C0B0BD59E2CA4BDB556D61A5", 16)
	secp224k1.BitSize = 224
}

// S224 returns a BitCurve which implements secp224k1 (see SEC 2 section 2.6.1)
func S224() *BitCurve {
	initonce.Do(initAll)
	return secp224k1
}
