package chromath

const (
	CIEKappa = 24389.0/27.0
	CIEEps = 216/24389.0
)

type Point [3]float64
type XYZ Point

type YZMatrix Matrix

func (p XYZ) X() float64 {
	return p[0]
}

func (p XYZ) Y() float64 {
	return p[1]
}

func (p XYZ) Z() float64 {
	return p[2]
}

// Transform applies an arbitrary transformation matrix to an XYZ point
// This is most useful with adaptation matrices generated from (Adaptation) Transform
func (p XYZ) Transform(tm *Matrix) XYZ {
	return XYZ(tm.Mul3x1(Point(p)))
}

type Lab Point

func (p Lab) L() float64 {
	return p[0]
}

func (p Lab) A() float64 {
	return p[1]
}

func (p Lab) B() float64 {
	return p[2]
}

type LCh Point

func (p LCh) L() float64 {
	return p[0]
}

func (p LCh) C() float64 {
	return p[1]
}

func (p LCh) H() float64 {
	return p[2]
}

type Luv Point

func (p Luv) L() float64 {
	return p[0]
}

func (p Luv) U() float64 {
	return p[1]
}

func (p Luv) V() float64 {
	return p[2]
}

type LChuv Point

func (p LChuv) L() float64 {
	return p[0]
}

func (p LChuv) C() float64 {
	return p[1]
}

func (p LChuv) H() float64 {
	return p[2]
}

type RGB Point

func (p RGB) R() float64 {
	return p[0]
}

func (p RGB) G() float64 {
	return p[1]
}

func (p RGB) B() float64 {
	return p[2]
}

// XyYPrimary specifies the three additive primaries in CIE XyY
// The first index is for the primary (ie 0 is typically R in RGB) and the
// second index is an xyY element
type XyYPrimary struct {
	Xr, Yr, Xg, Yg, Xb, Yb float64
}

type Gamma float64

type Compander interface {
	Init(*RGBSpace) Compander
	Compand(Point) Point
	Linearize(Point) Point
}

type Scaler interface {
	Init(*RGBSpace) Scaler
	Scale(Point) Point
	ScaleInv(Point) Point
}

// XYZTransform calculates a transform matrix for the primaries to XYZ
func (p XyYPrimary) XYZTransform(illuminantRef XYZ) Matrix {
	m := Matrix{
		p.Xr/p.Yr, p.Xg/p.Yg, p.Xb/p.Yb,
		1.0, 1.0, 1.0,
		(1.0-p.Xr-p.Yr)/p.Yr, (1.0-p.Xg-p.Yg)/p.Yg, (1.0-p.Xb-p.Yb)/p.Yb,
	}
	s := m.Inv().Transpose().Mul3x1(Point(illuminantRef))

	return Matrix{
		s[0]*m[0], s[1]*m[1], s[2]*m[2],
		s[0]*m[3], s[1]*m[4], s[2]*m[5],
		s[0]*m[6], s[1]*m[7], s[2]*m[8],
	}.Transpose()
}

type Adaptation Matrix

func (a Adaptation) Transform(sourceIlluminant XYZ, targetIlluminant XYZ) Matrix {
	ma := Matrix(a)
	mai := ma.Inv()
	sc := ma.Mul3x1(Point(sourceIlluminant))
	dc := ma.Mul3x1(Point(targetIlluminant))

	sm := Matrix{dc[0]/sc[0], 0, 0, 0, dc[1]/sc[1], 0, 0, 0, dc[2]/sc[2]}
	return mai.Mul3(sm).Mul3(ma)
}

func sqr(v float64) float64 {
	return v * v
}
