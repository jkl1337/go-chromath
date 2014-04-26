// Package deltae provides implementations of common CIE standard color difference (ΔE)
// calculations.
package deltae

import (
	. "github.com/jkl1337/go-chromath"
	"math"
)

const (
	pi = math.Pi
)

func sqr(v float64) float64 {
	return v * v
}

// KLCh represents the weighting parameters that are used for the CIEDE2000 color difference calculation.
type KLCh struct {
	KL, KC, Kh float64
}

// KLCHDefault is the most commonly used set of weighting parameters for CIEDE2000
var KLChDefault = KLCh{1, 1, 1}

// CIE2000 computes the CIEDE2000 delta-E for two L*a*b* space color coordinates
// klch is for configuring the weighting factors, but this almost always should be KLCHDefault
// Note that this implementation will exhibit slightly different behavior around the discontinuities
// of the function (these are grey colors) compared to Java and most C runtimes. The golang atan
// function has different accuracy characteristics compared to most Unix platforms and Java Strict math
func CIE2000(std Lab, sample Lab, klch *KLCh) float64 {
	lBarPrime := (std.L() + sample.L()) * 0.5
	c1 := math.Sqrt(std.A()*std.A() + std.B()*std.B())
	c2 := math.Sqrt(sample.A()*sample.A() + sample.B()*sample.B())
	cBar := (c1 + c2) * 0.5

	cBar7 := cBar * cBar * cBar
	cBar7 *= cBar7 * cBar
	g := 0.5 * (1.0 - math.Sqrt(cBar7/(cBar7+6103515625.0))) // 25**7

	a1Prime := (1.0 + g) * std.A()
	a2Prime := (1.0 + g) * sample.A()

	c1Prime := math.Sqrt(a1Prime*a1Prime + std.B()*std.B())
	c2Prime := math.Sqrt(a2Prime*a2Prime + sample.B()*sample.B())

	cBarPrime := (c1Prime + c2Prime) * 0.5

	h1Prime := math.Atan2(std.B(), a1Prime)
	if h1Prime < 0 {
		h1Prime += 2 * pi
	}
	h2Prime := math.Atan2(sample.B(), a2Prime)
	if h2Prime < 0 {
		h2Prime += 2 * pi
	}

	hBarPrime := (h1Prime + h2Prime) * 0.5
	dhPrime := h2Prime - h1Prime
	if math.Abs(dhPrime) > pi {
		hBarPrime += pi
		if h2Prime <= h1Prime {
			dhPrime += 2 * pi
		} else {
			dhPrime -= 2 * pi
		}
	}

	t := 1.0 -
		0.17*math.Cos(hBarPrime-pi/6) +
		0.24*math.Cos(2.0*hBarPrime) +
		0.32*math.Cos(3.0*hBarPrime+pi/30) -
		0.20*math.Cos(4.0*hBarPrime-63.0*pi/180)

	dLPrime := sample.L() - std.L()
	dCPrime := c2Prime - c1Prime
	dHPrime := 2.0 * math.Sqrt(c1Prime*c2Prime) * math.Sin(dhPrime/2.0)

	lBarPrimeM50Sqr := lBarPrime - 50.0
	lBarPrimeM50Sqr *= lBarPrimeM50Sqr
	sL := 1.0 + (0.015*lBarPrimeM50Sqr)/math.Sqrt(20.0+lBarPrimeM50Sqr)
	sC := 1.0 + 0.045*cBarPrime
	sH := 1.0 + 0.015*cBarPrime*t

	hBarPrimeM := (180/pi*hBarPrime - 275.0) / 25.0
	dTheta := pi / 6 * math.Exp(-hBarPrimeM*hBarPrimeM)
	cBarPrime7 := cBarPrime * cBarPrime * cBarPrime
	cBarPrime7 *= cBarPrime7 * cBarPrime
	rC := math.Sqrt(cBarPrime7 / (cBarPrime7 + 6103515625.0))
	rT := -2.0 * rC * math.Sin(2.0*dTheta)

	return math.Sqrt(
		sqr(dLPrime/(klch.KL*sL)) +
			sqr(dCPrime/(klch.KC*sC)) +
			sqr(dHPrime/(klch.Kh*sH)) +
			(dCPrime/(klch.KC*sC))*(dHPrime/(klch.Kh*sH))*rT)

}

// CIE76 computes the CIE76 color difference. This is just Euclidean distance in Lab space, and therefore quite fast, though it exhibits perceptual uniformity issues especially in the blue and desaturated regions.
func CIE76(std Lab, sample Lab) float64 {
	return math.Sqrt(sqr(std.L()-sample.L()) + sqr(std.A()-sample.A()) + sqr(std.B()-sample.B()))
}

// KLCh94 is a struct for weighting factors for CIE94 ΔE calculation.
type KLCh94 struct {
	KL, KC, Kh, K1, K2 float64
}

// KLCH94GraphicArts are the weighting factors for CIE94 used for most uses except textiles.
var KLCH94GraphicArts = KLCh94{1, 1, 1, 0.045, 0.015}

// KLCH94Textiles are the weighting factors for CIE94 used for textiles.
var KLCH94Textiles = KLCh94{2, 1, 1, 0.048, 0.014}

// DeltaECIE94 computes the CIE94 color difference of two L*a*b* colors.
// This is a distance calculation with the addition of weighting factors specified by klch.
func CIE94(std Lab, sample Lab, klch *KLCh94) float64 {
	dLsq := sqr(std.L() - sample.L())
	c1 := math.Sqrt(sqr(std.A()) + sqr(std.B()))
	c2 := math.Sqrt(sqr(sample.A()) + sqr(sample.B()))

	dCsq := sqr(c1 - c2)
	dHsq := sqr(std.A()-sample.A()) + sqr(std.B()-sample.B()) - dCsq

	sC := 1.0 + klch.K1*c1
	sH := 1.0 + klch.K2*c1

	return math.Sqrt(
		dLsq/sqr(klch.KL) +
			dCsq/sqr(klch.KC*sC) +
			dHsq/sqr(klch.Kh*sH))
}
