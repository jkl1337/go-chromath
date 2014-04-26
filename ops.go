package chromath

import (
	"math"
	"sync"
)

func clip(v float64) float64 {
	if v < 0.0 {
		v = 0.0
	} else if v > 1.0 {
		v = 1.0
	}
	return v
}

type scaler8bClamping struct{}

// Scaler8bClamping is a simple mapping of values [0,255] to [0,1], clamping all out of bounds values
var Scaler8bClamping scaler8bClamping

func (*scaler8bClamping) Init(ws *RGBSpace) Scaler { return &Scaler8bClamping }

func (*scaler8bClamping) Scale(p Point) Point {
	p[0] *= (1/255.0)
	p[1] *= (1/255.0)
	p[2] *= (1/255.0)
	return p
}

func (*scaler8bClamping) ScaleInv(p Point) Point {
	p[0] = math.Floor(0.5 + clip(p[0])*255)
	p[1] = math.Floor(0.5 + clip(p[1])*255)
	p[2] = math.Floor(0.5 + clip(p[2])*255)
	return p
}

type scaler16bClamping struct{}

// Scaler16bClamping is a simple mapping of values [0,2^16-1] to [0,1], clamping all out of bounds values. Typically used for 16-bit RGB
var Scaler16bClamping scaler16bClamping

func (*scaler16bClamping) Init(ws *RGBSpace) Scaler { return &Scaler16bClamping }

func (*scaler16bClamping) Scale(p Point) Point {
	p[0] *= (1/65535.0)
	p[1] *= (1/65535.0)
	p[2] *= (1/65535.0)
	return p
}

func (*scaler16bClamping) ScaleInv(p Point) Point {
	p[0] = math.Floor(0.5 + clip(p[0])*65535)
	p[1] = math.Floor(0.5 + clip(p[1])*65535)
	p[2] = math.Floor(0.5 + clip(p[2])*65535)
	return p
}

type gammaCompander struct {
	gamma, igamma float64
}
// GammaCompander is used for RGB spaces with a simple power (γ) function
var GammaCompander gammaCompander

func (*gammaCompander) Init(ws *RGBSpace) Compander {
	return &gammaCompander{
		gamma: float64(ws.Gamma),
		igamma: float64(1/ws.Gamma),
	}
}

func (c *gammaCompander) Compand(p Point) Point {
	igamma := c.igamma

	if p[0] < 0.0 {
		p[0] = -math.Exp(igamma * math.Log(-p[0]))
	} else {
		p[0] = math.Exp(igamma * math.Log(p[0]))
	}

	if p[1] < 0.0 {
		p[1] = -math.Exp(igamma * math.Log(-p[1]))
	} else {
		p[1] = math.Exp(igamma * math.Log(p[1]))
	}

	if p[2] < 0.0 {
		p[2] = -math.Exp(igamma * math.Log(-p[2]))
	} else {
		p[2] = math.Exp(igamma * math.Log(p[2]))
	}
	return p
}

func (c *gammaCompander) Linearize(p Point) Point {
	gamma := c.gamma

	if p[0] < 0.0 {
		p[0] = -math.Exp(gamma * math.Log(-p[0]))
	} else {
		p[0] = math.Exp(gamma * math.Log(p[0]))
	}

	if p[1] < 0.0 {
		p[1] = -math.Exp(gamma * math.Log(-p[1]))
	} else {
		p[1] = math.Exp(gamma * math.Log(p[1]))
	}

	if p[2] < 0.0 {
		p[2] = -math.Exp(gamma * math.Log(-p[2]))
	} else {
		p[2] = math.Exp(gamma * math.Log(p[2]))
	}
	return p
}

type sRGBCompander struct {}
// SRGBCompander provides the companding γ function as specified for sRGB.
var SRGBCompander sRGBCompander

func (*sRGBCompander) Init(ws *RGBSpace) Compander { return &SRGBCompander }

func (*sRGBCompander) Compand(p Point) Point {
	sign := 1.0
	if p[0] < 0.0 {
		sign = -1.0
		p[0] = -p[0]
	}
	if p[0] > 0.0031308 {
		p[0] = 1.055*math.Exp((1/2.4)*math.Log(p[0])) - 0.055
	} else {
		p[0] = p[0] * 12.92
	}
	p[0] *= sign

	sign = 1.0
	if p[1] < 0.0 {
		sign = -1.0
		p[1] = -p[1]
	}
	if p[1] > 0.0031308 {
		p[1] = 1.055*math.Exp((1/2.4)*math.Log(p[1])) - 0.055
	} else {
		p[1] = p[1] * 12.92
	}
	p[1] *= sign

	sign = 1.0
	if p[2] < 0.0 {
		sign = -1.0
		p[2] = -p[2]
	}
	if p[2] > 0.0031308 {
		p[2] = 1.055*math.Exp((1/2.4)*math.Log(p[2])) - 0.055
	} else {
		p[2] = p[2] * 12.92
	}
	p[2] *= sign

	return p
}

func (*sRGBCompander) Linearize(p Point) Point {
	sign := 1.0
	if p[0] < 0.0 {
		sign = -1.0
		p[0] = -p[0]
	}
	if p[0] <= 0.04045 {
		p[0] = p[0] / 12.92
	} else {
		p[0] = math.Exp(2.4 * math.Log((p[0]+0.055)/1.055))
	}
	p[0] *= sign

	sign = 1.0
	if p[1] < 0.0 {
		sign = -1.0
		p[1] = -p[1]
	}
	if p[1] <= 0.04045 {
		p[1] = p[1] / 12.92
	} else {
		p[1] = math.Exp(2.4 * math.Log((p[1]+0.055)/1.055))
	}
	p[1] *= sign

	sign = 1.0
	if p[2] < 0.0 {
		sign = -1.0
		p[2] = -p[2]
	}
	if p[2] <= 0.04045 {
		p[2] = p[2] / 12.92
	} else {
		p[2] = math.Exp(2.4 * math.Log((p[2]+0.055)/1.055))
	}
	p[2] *= sign

	return p
}

var sRGBFastCompanderOnce sync.Once
var sRGBFastCompanderTable []float64

type sRGBFastCompander struct {}
// SRGBFastCompander provides both a compander and scaler for 8-bit sRGB using a lookup table for linearization
// This compander should not be used with a scaler since it provides both scaling and the γ function.
var SRGBFastCompander sRGBFastCompander

func sRGBFastCompanderInit() {
	sRGBFastCompanderTable = make([]float64, 256)

	for i := 0; i < 256; i++ {
		v := float64(i)/255.0
		if v <= 0.04045 {
			v = v / 12.92
		} else {
			v = math.Exp(2.4 * math.Log((v+0.055)/1.055))
		}
		sRGBFastCompanderTable[i] = v
	}
}

func (*sRGBFastCompander) Init(ws *RGBSpace) Compander {
	sRGBFastCompanderOnce.Do(sRGBFastCompanderInit)
	return &SRGBFastCompander
}

func (*sRGBFastCompander) Compand(p Point) Point {
	return Scaler8bClamping.ScaleInv(SRGBCompander.Compand(p))
}

func vInRange(n float64) float64 {
	i := int(n)
	if i < 0 {
		i = 0
	} else if i > 255 {
		i = 255
	}
	return sRGBFastCompanderTable[i]
}

func (*sRGBFastCompander) Linearize(p Point) Point {
	return Point{vInRange(p[0]), vInRange(p[1]), vInRange(p[2])}
}

type lstarCompander struct {}
// LstarCompander is a compander used for L* gamma response used in working spaces such
// as ECI RGB and in alternative RGB calibrated workflows.
var LstarCompander lstarCompander

func (*lstarCompander) Init(ws *RGBSpace) Compander { return &LstarCompander }

func (*lstarCompander) Compand(p Point) Point {
	l := p[0]
	if l < 0.0 { l = -l }
	if l <= 216.0/24389.0 {
		l = l * 24389.0 / 2700.0
	} else {
		l = 1.16 * math.Cbrt(l) - 0.16
	}
	if p[0] < 0.0 {
		l = -l
	}
	p[0] = l

	l = p[1]
	if l < 0.0 { l = -l }
	if l <= 216.0/24389.0 {
		l = l * 24389.0 / 2700.0
	} else {
		l = 1.16 * math.Cbrt(l) - 0.16
	}
	if p[1] < 0.0 {
		l = -l
	}
	p[1] = l

	l = p[2]
	if l < 0.0 { l = -l }
	if l <= 216.0/24389.0 {
		l = l * 24389.0 / 2700.0
	} else {
		l = 1.16 * math.Cbrt(l) - 0.16
	}
	if p[2] < 0.0 {
		l = -l
	}
	p[2] = l
	return p
}

func (*lstarCompander) Linearize(p Point) Point {
	c := p[0]
	if c < 0.0 { c = -c }
	if c <= 0.08 {
		c = 2700.0 * c / 24389.0
	} else {
		c = (((1000000.0 * c + 480000.0) * c + 76800.0) * c + 4096.0) / 1560896.0
	}
	if p[0] < 0.0 {
		c = -c
	}
	p[0] = c

	c = p[1]
	if c < 0.0 { c = -c }
	if c <= 0.08 {
		c = 2700.0 * c / 24389.0
	} else {
		c = (((1000000.0 * c + 480000.0) * c + 76800.0) * c + 4096.0) / 1560896.0
	}
	if p[1] < 0.0 {
		c = -c
	}
	p[1] = c

	c = p[2]
	if c < 0.0 { c = -c }
	if c <= 0.08 {
		c = 2700.0 * c / 24389.0
	} else {
		c = (((1000000.0 * c + 480000.0) * c + 76800.0) * c + 4096.0) / 1560896.0
	}
	if p[2] < 0.0 {
		c = -c
	}
	p[2] = c

	return p
}

type bT2020Compander struct {}
// BT2020Compander provides a compander matching the transfer function speicied by the BT.2020 (HDTV) recommendation
// It should be usable for both 12-bit and 10-bit simulations
var BT2020Compander bT2020Compander

func (*bT2020Compander) Init(ws *RGBSpace) Compander { return &BT2020Compander }

func (*bT2020Compander) Compand(p Point) Point {
	if p[0] < 0.0181 {
		p[0] = 4.5 * p[0]
	} else {
		p[0] = 1.0993 * math.Exp(0.45 * math.Log(p[0])) - 0.099
	}
	if p[1] < 0.0181 {
		p[1] = 4.5 * p[1]
	} else {
		p[1] = 1.0993 * math.Exp(0.45 * math.Log(p[1])) - 0.099
	}
	if p[2] < 0.0181 {
		p[2] = 4.5 * p[2]
	} else {
		p[2] = 1.0993 * math.Exp(0.45 * math.Log(p[2])) - 0.099
	}
	return p
}

func (*bT2020Compander) Linearize(p Point) Point {
	if p[0] < 0.08145 {
		p[0] /= 4.5
	} else {
		p[0] = math.Exp(1/0.45 * math.Log((p[0] + 0.099)/1.099))
	}
	if p[1] < 0.08145 {
		p[1] /= 4.5
	} else {
		p[1] = math.Exp(1/0.45 * math.Log((p[1] + 0.099)/1.099))
	}
	if p[2] < 0.08145 {
		p[2] /= 4.5
	} else {
		p[2] = math.Exp(1/0.45 * math.Log((p[2] + 0.099)/1.099))
	}
	return p
}
