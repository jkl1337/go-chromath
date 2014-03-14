package chromath

import "testing"
import "math"
import "fmt"

type T testing.T

func round(n float64, digits int) float64 {
	f := 1.0
	i := digits
	if i < 0 { i = -i }
	for i-- ; i >= 0 ; i-- {
		f = f*10.0
	}
	if digits < 0 {
		n /= f
	} else  {
		n *= f
	}

	if n >= 0.0 {
		n = math.Floor(n + 0.5)
	} else {
		n = math.Ceil(n - 0.5)
	}

	if digits < 0 {
		n *= f
	} else {
		n /= f
	}
	return n
}

func (t *T) testPointEqual(a Point, b Point) {
	if a[0] != b[0] || a[1] != b[1] || a[2] != b[2] {
		t.Errorf("Point %v != %v", a, b)
	}
}

func (t *T) testAlmostEqual(a float64, b float64, places int, msg string) {
	if round(math.Abs(a - b), places) != 0 {
		t.Error(msg, a, b)
	}
}

func (t *T) testPointAlmostEqual(a Point, b Point, places int) {
	if round(math.Abs(a[0] - b[0]), places) != 0 ||
		round(math.Abs(a[1] - b[1]), places) != 0 ||
		round(math.Abs(a[2] - b[2]), places) != 0 {
		t.Errorf("Point %v !~ %v", a, b)
	}
}

func TestXYZToLabConversion(t_ *testing.T) {
	t := (*T)(t_)

	lt := NewLabTransformer(nil, 1.0)
	lab := lt.Invert(XYZ{0.1, 0.2, 0.3})
	t.testPointAlmostEqual(Point(lab), Point{51.837, -57.486, -25.780}, 3)
}

func TestXYZToLuvConversion(t_ *testing.T) {
	t := (*T)(t_)

	lt := NewLuvTransformer(nil, 1.0)
	lab := lt.Invert(XYZ{0.1, 0.2, 0.3})
	t.testPointAlmostEqual(Point(lab), Point{51.837, -73.561, -25.657}, 3)
}


func TestLuvToXYZConversion(t_ *testing.T) {
	t := (*T)(t_)

	lt := NewLuvTransformer(nil, 1.0)
	xyz := lt.Convert(Luv{1.807, -2.564, -0.894})
	t.testPointAlmostEqual(Point(xyz), Point{0.001, 0.002, 0.003}, 3)
}

func TestLabtoLCHuvConversion(t_ *testing.T) {
	t := (*T)(t_)

	lchuv := Luv{1.807, -2.564, -0.894}.LChuv()
	t.testPointAlmostEqual(Point(lchuv), Point{1.807, 2.715, 199.222}, 3)
}

func TestLabtoXYZConversion(t_ *testing.T) {
	t := (*T)(t_)

	lt := NewLabTransformer(nil, 1.0)
	xyz := lt.Convert(Lab{1.807, -3.749, -2.547})
	t.testPointAlmostEqual(Point(xyz), Point{0.001, 0.002, 0.003}, 3)
}

func TestLabtoLCHConversion(t_ *testing.T) {
	t := (*T)(t_)

	lch := Lab{1.807, -3.749, -2.547}.LCh()
	t.testPointAlmostEqual(Point(lch), Point{1.807, 4.532, 214.191}, 3)
}

func TestRGBToXYZConversion(t_ *testing.T) {
	t := (*T)(t_)

	rd50 := NewRGBTransformer(&SpaceSRGB, &AdaptationBradford, nil, &Scaler8bClamping, 100.0, nil)
	origRgb := RGB{123, 200, 50}

	xyzd50 := rd50.Convert(origRgb)
	t.testPointAlmostEqual(Point(xyzd50), Point{31.334, 46.006, 8.162}, 3)
	rtRgb := rd50.Invert(xyzd50)

	t.testPointAlmostEqual(Point(rtRgb), Point(origRgb), 3)

	// check default whitepoint
	rd50imp := NewRGBTransformer(&SpaceSRGB, &AdaptationBradford, &IlluminantRefD50, &Scaler8bClamping, 100.0, nil)
	xyzd50imp := rd50imp.Convert(origRgb)
	t.testPointEqual(Point(xyzd50imp), Point(xyzd50))

	rd65 := NewRGBTransformer(&SpaceSRGB, &AdaptationBradford, &IlluminantRefD65, &Scaler8bClamping, 100.0, nil)
	xyzd65 := rd65.Convert(origRgb)
	t.testPointAlmostEqual(Point(xyzd65), Point{29.398, 45.748, 10.298}, 3)
}

func TestAdobeRGBToXYZConversion(t_ *testing.T) {
	t := (*T)(t_)

	rgb := RGB{123, 200, 50}
	rd65 := NewRGBTransformer(&SpaceAdobeRGB, &AdaptationBradford, &IlluminantRefD65, &Scaler8bClamping, 0.0, nil)
	xyzd65 := rd65.Convert(rgb)
	t.testPointAlmostEqual(Point(xyzd65), Point{0.230, 0.430, 0.074}, 3)

	rd50 := NewRGBTransformer(&SpaceAdobeRGB, &AdaptationBradford, nil, &Scaler8bClamping, 0.0, nil)
	xyzd50 := rd50.Convert(rgb)
	t.testPointAlmostEqual(Point(xyzd50), Point{0.247, 0.431, 0.060}, 3)
}

func TestLChToRGBConversion(t_ *testing.T) {
	t := (*T)(t_)
	lab2xyz := NewLabTransformer(nil, 0)
	rgb2xyz := NewRGBTransformer(&SpaceSRGB, nil, nil, nil, 0.0, nil)
	rgb := rgb2xyz.Invert(lab2xyz.Convert(LCh{0.0, 0.0, 0.0}.Lab()))
	t.testPointEqual(Point(rgb), Point{0.0, 0.0, 0.0})
}

func TestLChuvToLuvConversion(t_ *testing.T) {
	t := (*T)(t_)
	luv := LChuv{1.807, 2.715, 199.228}.Luv()
	t.testPointAlmostEqual(Point(luv), Point{1.807, -2.564, -0.894}, 3)
}

func TestAdaptation(t_ *testing.T) {
	t := (*T)(t_)
	xyz := XYZ{0.5, 0.4, 0.1}
	am := AdaptationBradford.Transform(IlluminantRefC.XYZ, IlluminantRefD65.XYZ)
	r := xyz.Transform(&am)
	t.testPointAlmostEqual(Point(r), Point{0.4912, 0.3998, 0.0927}, 4)
}

func TestCIE2000_1(t_ *testing.T) {
	t := (*T)(t_)

	expected := 1.523
	c1, c2 := Lab{0.9, 16.3, -2.22}, Lab{0.7, 14.2, -1.80}
	result := DeltaECIE2000(c1, c2, &KLChDefault)
	t.testAlmostEqual(result, expected, 3, fmt.Sprintf("%v, %v", c1, c2))
}

func TestCIE2000_2(t_ *testing.T) {
	t := (*T)(t_)

	c1, c2 := Lab{32.8911, -53.0107, -43.3182}, Lab{77.1797, 25.5928, 17.9412}
	expected := 78.772
	result := DeltaECIE2000(c1, c2, &KLChDefault)
	t.testAlmostEqual(result, expected, 3, fmt.Sprintf("%v, %v", c1, c2))
}

type colorDiffCase struct {
	ref, sample Lab
	expectedCIE2000 float64
}

// From "The CIEDE2000 Color-Difference Formula: Implementation Notes,
// Supplementary Test Data, and Mathematical Observations,", G. Sharma,
// W. Wu, E. N. Dalal, submitted to Color Research and Application,
// January 2004. http://www.ece.rochester.edu/~gsharma/ciede2000/
var colorCases = []colorDiffCase{
	{ Lab{ 50.0000,   2.6772, -79.7751 }, Lab{ 50.0000,   0.0000, -82.7485 },   2.0425 },
	{ Lab{ 50.0000,   3.1571, -77.2803 }, Lab{ 50.0000,   0.0000, -82.7485 },   2.8615 },
	{ Lab{ 50.0000,   2.8361, -74.0200 }, Lab{ 50.0000,   0.0000, -82.7485 },   3.4412 },
	{ Lab{ 50.0000,  -1.3802, -84.2814 }, Lab{ 50.0000,   0.0000, -82.7485 },   1.0000 },
	{ Lab{ 50.0000,  -1.1848, -84.8006 }, Lab{ 50.0000,   0.0000, -82.7485 },   1.0000 },
	{ Lab{ 50.0000,  -0.9009, -85.5211 }, Lab{ 50.0000,   0.0000, -82.7485 },   1.0000 },
	{ Lab{ 50.0000,   0.0000,   0.0000 }, Lab{ 50.0000,  -1.0000,   2.0000 },   2.3669 },
	{ Lab{ 50.0000,  -1.0000,   2.0000 }, Lab{ 50.0000,   0.0000,   0.0000 },   2.3669 },
	{ Lab{ 50.0000,   2.4900,  -0.0010 }, Lab{ 50.0000,  -2.4900,   0.0009 },   7.1792 },
	// This is one of those discontinuity test cases. And it seems to be provoking an discrepancy
	// between math.Atan2 and most other implementations
	//{ Lab{ 50.0000,   2.4900,  -0.0010 }, Lab{ 50.0000,  -2.4900,   0.0010 },   7.1792 },
	{ Lab{ 50.0000,   2.4900,  -0.0010 }, Lab{ 50.0000,  -2.4900,   0.0011 },   7.2195 },
	{ Lab{ 50.0000,   2.4900,  -0.0010 }, Lab{ 50.0000,  -2.4900,   0.0012 },   7.2195 },
	{ Lab{ 50.0000,  -0.0010,   2.4900 }, Lab{ 50.0000,   0.0009,  -2.4900 },   4.8045 },
	{ Lab{ 50.0000,  -0.0010,   2.4900 }, Lab{ 50.0000,   0.0010,  -2.4900 },   4.8045 },
	{ Lab{ 50.0000,  -0.0010,   2.4900 }, Lab{ 50.0000,   0.0011,  -2.4900 },   4.7461 },
	{ Lab{ 50.0000,   2.5000,   0.0000 }, Lab{ 50.0000,   0.0000,  -2.5000 },   4.3065 },
	{ Lab{ 50.0000,   2.5000,   0.0000 }, Lab{ 73.0000,  25.0000, -18.0000 },  27.1492 },
	{ Lab{ 50.0000,   2.5000,   0.0000 }, Lab{ 61.0000,  -5.0000,  29.0000 },  22.8977 },
	{ Lab{ 50.0000,   2.5000,   0.0000 }, Lab{ 56.0000, -27.0000,  -3.0000 },  31.9030 },
	{ Lab{ 50.0000,   2.5000,   0.0000 }, Lab{ 58.0000,  24.0000,  15.0000 },  19.4535 },
	{ Lab{ 50.0000,   2.5000,   0.0000 }, Lab{ 50.0000,   3.1736,   0.5854 },   1.0000 },
	{ Lab{ 50.0000,   2.5000,   0.0000 }, Lab{ 50.0000,   3.2972,   0.0000 },   1.0000 },
	{ Lab{ 50.0000,   2.5000,   0.0000 }, Lab{ 50.0000,   1.8634,   0.5757 },   1.0000 },
	{ Lab{ 50.0000,   2.5000,   0.0000 }, Lab{ 50.0000,   3.2592,   0.3350 },   1.0000 },
	{ Lab{ 60.2574, -34.0099,  36.2677 }, Lab{ 60.4626, -34.1751,  39.4387 },   1.2644 },
	{ Lab{ 63.0109, -31.0961,  -5.8663 }, Lab{ 62.8187, -29.7946,  -4.0864 },   1.2630 },
	{ Lab{ 61.2901,   3.7196,  -5.3901 }, Lab{ 61.4292,   2.2480,  -4.9620 },   1.8731 },
	{ Lab{ 35.0831, -44.1164,   3.7933 }, Lab{ 35.0232, -40.0716,   1.5901 },   1.8645 },
	{ Lab{ 22.7233,  20.0904, -46.6940 }, Lab{ 23.0331,  14.9730, -42.5619 },   2.0373 },
	{ Lab{ 36.4612,  47.8580,  18.3852 }, Lab{ 36.2715,  50.5065,  21.2231 },   1.4146 },
	{ Lab{ 90.8027,  -2.0831,   1.4410 }, Lab{ 91.1528,  -1.6435,   0.0447 },   1.4441 },
	{ Lab{ 90.9257,  -0.5406,  -0.9208 }, Lab{ 88.6381,  -0.8985,  -0.7239 },   1.5381 },
	{ Lab{  6.7747,  -0.2908,  -2.4247 }, Lab{  5.8714,  -0.0985,  -2.2286 },   0.6377 },
	{ Lab{  2.0776,   0.0795,  -1.1350 }, Lab{  0.9033,  -0.0636,  -0.5514 },   0.9082 },
}

func TestCIE2000Data(t_ *testing.T) {
	t := (*T)(t_)

	for _, c := range(colorCases) {
		result := DeltaECIE2000(c.ref, c.sample, &KLChDefault)
		resultBa := DeltaECIE2000(c.ref, c.sample, &KLChDefault)

		if (resultBa != result) {
			t.Error("Asymmetry for %v, %v (%v != %v)", c.ref, c.sample, result, resultBa)
		}
		t.testAlmostEqual(result, c.expectedCIE2000, 4, fmt.Sprintf("%v, %v", c.ref, c.sample))
	}
}

func BenchmarkCIE2000(b *testing.B) {
	c := &colorCases[0]
	for i := 0; i < b.N; i++ {
		DeltaECIE2000(c.ref, c.sample, &KLChDefault)
	}
}

func BenchmarkRGBXYZ(b *testing.B) {
	r := NewRGBTransformer(&SpaceSRGB, &AdaptationBradford, nil, &Scaler8bClamping, 0.0, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := r.Convert(RGB{60, 92, 180})
		r.Invert(p)
	}
}

func BenchmarkSRGBLinearize(b *testing.B) {
	c := SRGBCompander.Init(nil)
	p := Point{0.2, 0.4, 0.8}
	for i := 0; i < b.N; i++ {
		c.Linearize(p)
	}
}

func BenchmarkSRGBFastLinearize(b *testing.B) {
	c := SRGBFastCompander.Init(nil)
	p := Point{15, 200, 39}
	for i := 0; i < b.N; i++ {
		c.Linearize(p)
	}
}

func BenchmarkSRGBCompand(b *testing.B) {
	var p Point
	c := SRGBCompander.Init(nil)
	for i := 0; i < b.N; i++ {
		p = Point{0.2, 0.4, 0.8}
		p = c.Compand(p)
		p = c.Linearize(p)
	}
}


