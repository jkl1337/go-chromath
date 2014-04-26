package chromath

import (
	"math"
)

// LabTransformer performs transform to CIELAB from XYZ and vice-versa
type LabTransformer struct {
	refWp XYZ
}

// NewLabTransformer creates a LAB ⇔ XYZ transformer
func NewLabTransformer(refIlluminant *IlluminantRef) *LabTransformer {
	if refIlluminant == nil {
		refIlluminant = &IlluminantRefD50
	}
	return &LabTransformer{refIlluminant.XYZ}
}

// Convert converts a LAB point to XYZ
func (t *LabTransformer) Convert(p Lab) XYZ {
	fy := (p.L() + 16.0) / 116.0
	fx := 0.002 * p.A() + fy
	fz := fy - 0.005 * p.B()

	fx3 := fx * fx * fx
	fz3 := fz * fz * fz

	var xr, yr, zr float64
	if fx3 > CIEEps {
		xr = fx3
	} else {
		xr = (116.0 * fx - 16.0) / CIEKappa
	}

	if p.L() > CIEKappa*CIEEps {
		yr = fy * fy * fy
	} else {
		yr = p.L() / CIEKappa
	}

	if fz3 > CIEEps {
		zr = fz3
	} else {
		zr = (116.0 * fz - 16.0) / CIEKappa
	}
	return XYZ{xr * t.refWp.X(), yr * t.refWp.Y(), zr * t.refWp.Z()}
}

// Invert converts an XYZ point to LAB
func (t *LabTransformer) Invert(p XYZ) Lab {
	xr, yr, zr := p.X()/t.refWp.X(), p.Y()/t.refWp.Y(), p.Z()/t.refWp.Z()

	var fx, fy, fz float64
	if xr > CIEEps {
		fx = math.Cbrt(xr)
	} else {
		fx = (CIEKappa * xr + 16.0) / 116.0
	}
	if yr > CIEEps {
		fy = math.Cbrt(yr)
	} else {
		fy = (CIEKappa * yr + 16.0) / 116.0
	}
	if zr > CIEEps {
		fz = math.Cbrt(zr)
	} else {
		fz = (CIEKappa * zr + 16.0) / 116.0
	}

	return Lab{116 * fy - 16.0, 500.0 * (fx - fy), 200.0 * (fy - fz)}
}

// LCh returns the LAB point as an LCh(ab) point (Lab in cylindrical coordinates)
func (p Lab) LCh() LCh {
	c := math.Sqrt(sqr(p.A()) + sqr(p.B()))
	h := 180.0 * math.Atan2(p.B(), p.A()) / math.Pi
	if h < 0.0 {
		h += 360.0
	}
	p[1], p[2] = c, h
	return LCh(p)
}

// Lab returns the LCh(ab) point as LAB
func (p LCh) Lab() Lab {
	a := p.C() * math.Cos(p.H() * math.Pi / 180.0)
	b := p.C() * math.Sin(p.H() * math.Pi / 180.0)
	p[1], p[2] = a, b
	return Lab(p)
}

// LCh2LabTransformer is a transformer from LAB to LCh(ab)
type LCh2LabTransformer struct {}
var lCh2LabTransformerInst LCh2LabTransformer

// NewLCh2LabTransformer creates a CIELAB ⇔ CIELCh(ab) transformer
func NewLCh2LabTransformer() *LCh2LabTransformer {
	return &lCh2LabTransformerInst
}

func (t *LCh2LabTransformer) Convert(p LCh) Lab {
	return p.Lab()
}

func (t *LCh2LabTransformer) Invert(p Lab) LCh {
	return p.LCh()
}
