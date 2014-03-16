package chromath

import (
	"math"
)

type LuvTransformer struct {
	refWp XYZ
	u0, v0 float64
}

func NewLuvTransformer(refIlluminant *IlluminantRef) *LuvTransformer {
	if refIlluminant == nil {
		refIlluminant = &IlluminantRefD50
	}
	wp := refIlluminant.XYZ
	dn := wp.X() + 15.0 * wp.Y() + 3.0 * wp.Z()
	u0 := 4.0 * wp.X() / dn
	v0 := 9.0 * wp.Y() / dn
	return &LuvTransformer{wp, u0, v0}
}

func (t *LuvTransformer) Convert(p Luv) XYZ {
	var xX, xY float64

	if p.L() > CIEEps*CIEKappa {
		xY = (p.L() + 16.0) / 116.0
		xY = xY * xY * xY
	} else {
		xY = p.L() / CIEKappa
	}

	a := (((52.0 * p.L()) / (p.U() + 13.0 * p.L() * t.u0)) - 1.0) / 3.0
	d := xY * (((39.0 * p.L()) / (p.V() + 13.0 * p.L() * t.v0)) - 5.0)
	b := -5.0 * xY
	c := -1.0 / 3.0

	xX = (d - b) / (a - c)
	return XYZ{xX, xY, xX * a + b}
}

func (t *LuvTransformer) Invert(p XYZ) Luv {
	d := p.X() + 15.0 * p.Y() + 3.0 * p.Z()

	var l, up, vp float64

	if d > 0 {
		up = 4.0 * p.X() / d
		vp = 9.0 * p.Y() / d
	}
	yr := p.Y() / t.refWp.Y()

	if yr > CIEEps {
		l = 116.0 * math.Cbrt(yr) - 16.0
	} else {
		l = CIEKappa * yr
	}
	return Luv{l, 13.0*l*(up-t.u0), 13.0*l*(vp-t.v0)}
}

func (p Luv) LChuv() LChuv {
	c := math.Sqrt(sqr(p.U()) + sqr(p.V()))
	h := 180.0 * math.Atan2(p.V(), p.U()) / math.Pi
	if h < 0.0 {
		h += 360.0
	}
	p[1], p[2] = c, h
	return LChuv(p)
}

func (p LChuv) Luv() Luv {
	u := p.C() * math.Cos(p.H() * math.Pi / 180.0)
	v := p.C() * math.Sin(p.H() * math.Pi / 180.0)
	p[1], p[2] = u, v
	return Luv(p)
}

type LChuv2LuvTransformer struct {}
var lChuv2LuvTransformerInst LChuv2LuvTransformer

func NewLChuv2LuvTransformer() *LChuv2LuvTransformer {
	return &lChuv2LuvTransformerInst
}

func (t *LChuv2LuvTransformer) Convert(p LChuv) Luv {
	return p.Luv()
}

func (t *LChuv2LuvTransformer) Invert(p Luv) LChuv {
	return p.LChuv()
}
