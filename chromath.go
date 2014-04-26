// The chromath package provides math functions for common colorspace transformations including
// customizable RGB working space to XYZ transformations and commonly used CIE colorspaces including CIELAB, CIELUV, and their cylindrical coordinate counterparts.
// Addtionally it provides convenient and performant application of commonly used chromatic adaptation transforms.
// Most transformations are invoked on methods of transformer objects. A transformer receiver
// can therefore precompute certain aspects of the transform such as a combined chromatic adaptation and working space transform matrix.
package chromath

// CIEKappa and CIEEps are the CIE defined constants (κ and ε) used for Lab and Luv transforms
const (
	CIEKappa = 24389.0 / 27.0
	CIEEps   = 216 / 24389.0
)

// Point is a generic 3-tuple colorspace point allowing for generic operations on points regardless of color space
type Point [3]float64

// XYZ is a CIEXYZ color coordinate. These colors need to be associated with a reference illuminant, stored externally.
type XYZ Point

func (p XYZ) X() float64 {
	return p[0]
}

func (p XYZ) Y() float64 {
	return p[1]
}

func (p XYZ) Z() float64 {
	return p[2]
}

// Transform applies an arbitrary transformation matrix to an XYZ point.
// This is most useful with adaptation matrices generated from (Adaptation) Transform.
func (p XYZ) Transform(tm *Matrix) XYZ {
	return XYZ(tm.Mul3x1(Point(p)))
}

// Lab represents a color in CIELab space.
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

// LCh represents a color in the CIELCH space which is exactly CIELAB represented with cylindrical coordinates.
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

// Luv represents a color in the CIELUV space.
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

// LCHuv represents a color in the CIELCH(uv) space which is exactly CIELUV represented with cylindrical coordinates.
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

// Gamma holds a gamma calculation parameter for use in working space definitions and as a
// parameter to RGB companding functions
type Gamma float64

// A Compander provides a bidirectional transforms from linear RGB to working RGB. It is
// defined this way so that a compander can be specified statically in a structure.
type Compander interface {
	// Init initializes a compander for use. It must be called on the object
	// implementing the compander before Compand or Linearize is called. A runtime panic
	// may occur if this is not done.
	Init(*RGBSpace) Compander
	// Compand transforms a linear RGB point to working RGB
	Compand(Point) Point
	// Linearize transforms a working RGB point to linear RGB
	Linearize(Point) Point
}

// A Scaler provides a scaling that is used for points in the working RGB space before they are
// linearized. The typical case is to scale 8-bit or 16-bit integers to the range [0,1] with clamping
type Scaler interface {
	// Init initializes a scaler instance for use. It must be called on the object
	// implementing the scaler before Scale or ScaleInv is called. A runtime panic
	// may occur if this is not done.
	Init(*RGBSpace) Scaler
	// Scale 
	Scale(Point) Point
	ScaleInv(Point) Point
}

// RGBTransform calculates a transform matrix given the primaries for RGB⇒XYZ
func (p XyYPrimary) RGBTransform(illuminantRef XYZ) Matrix {
	m := Matrix{
		p.Xr / p.Yr, p.Xg / p.Yg, p.Xb / p.Yb,
		1.0, 1.0, 1.0,
		(1.0 - p.Xr - p.Yr) / p.Yr, (1.0 - p.Xg - p.Yg) / p.Yg, (1.0 - p.Xb - p.Yb) / p.Yb,
	}
	s := m.Inv().Transpose().Mul3x1(Point(illuminantRef))

	return Matrix{
		s[0] * m[0], s[1] * m[1], s[2] * m[2],
		s[0] * m[3], s[1] * m[4], s[2] * m[5],
		s[0] * m[6], s[1] * m[7], s[2] * m[8],
	}.Transpose()
}

// An Adaptation is a chromatic adaptation transform matrix that defines the transform of XYZ colors into a cone response domain. The three most common adaptation matrices are included in this package.
type Adaptation Matrix

var (
	// AdaptationBradford is the most commonly used chromatic adapatation transform and is used in PhotoShop.
	AdaptationBradford = Adaptation{
		0.8951, -0.7502, 0.03890,
		0.2664, 1.7135, -0.06850,
		-0.1614000, 0.03670, 1.02960,
	}

	AdaptationVonKries = Adaptation{
		0.400240, -0.22630, 0.000,
		0.707600, 1.1653200, 0.000,
		-0.0808100, 0.045700, 0.918220,
	}

	AdaptationXYZScaling = Adaptation{1, 0, 0, 0, 1, 0, 0, 0, 1}
)

func sqr(v float64) float64 {
	return v * v
}

// Transform calculates the chromatic adapation matrix for the receiver given a source illuminant and target illuminant
func (a Adaptation) Transform(sourceIlluminant XYZ, targetIlluminant XYZ) Matrix {
	ma := Matrix(a)
	mai := ma.Inv()
	sc := ma.Mul3x1(Point(sourceIlluminant))
	dc := ma.Mul3x1(Point(targetIlluminant))

	sm := Matrix{dc[0] / sc[0], 0, 0, 0, dc[1] / sc[1], 0, 0, 0, dc[2] / sc[2]}
	return mai.Mul3(sm).Mul3(ma)
}
