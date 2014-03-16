package chromath

import (
	"testing"
	"math"
	"github.com/jkl1337/go-chromath/testnum"
)

func testPointEqual(t *testing.T, a Point, b Point) {
	if a[0] != b[0] || a[1] != b[1] || a[2] != b[2] {
		t.Errorf("Point %v != %v", a, b)
	}
}

func testPointAlmostEqual(t *testing.T, a Point, b Point, places int) {
	if testnum.Round(math.Abs(a[0] - b[0]), places) != 0 ||
		testnum.Round(math.Abs(a[1] - b[1]), places) != 0 ||
		testnum.Round(math.Abs(a[2] - b[2]), places) != 0 {
		t.Errorf("Point %v !~ %v", a, b)
	}
}

func TestXYZToLabConversion(t *testing.T) {
	lt := NewLabTransformer(nil)
	lab := lt.Invert(XYZ{0.1, 0.2, 0.3})
	testPointAlmostEqual(t, Point(lab), Point{51.837, -57.486, -25.780}, 3)
}

func TestXYZToLuvConversion(t *testing.T) {
	lt := NewLuvTransformer(nil)
	lab := lt.Invert(XYZ{0.1, 0.2, 0.3})
	testPointAlmostEqual(t, Point(lab), Point{51.837, -73.561, -25.657}, 3)
}

func TestLuvToXYZConversion(t *testing.T) {
	lt := NewLuvTransformer(nil)
	xyz := lt.Convert(Luv{1.807, -2.564, -0.894})
	testPointAlmostEqual(t, Point(xyz), Point{0.001, 0.002, 0.003}, 3)
}

func TestLabtoLCHuvConversion(t *testing.T) {
	lchuv := Luv{1.807, -2.564, -0.894}.LChuv()
	testPointAlmostEqual(t, Point(lchuv), Point{1.807, 2.715, 199.222}, 3)
}

func TestLabtoXYZConversion(t *testing.T) {
	lt := NewLabTransformer(nil)
	xyz := lt.Convert(Lab{1.807, -3.749, -2.547})
	testPointAlmostEqual(t, Point(xyz), Point{0.001, 0.002, 0.003}, 3)
}

func TestLabtoLCHConversion(t *testing.T) {
	lch := Lab{1.807, -3.749, -2.547}.LCh()
	testPointAlmostEqual(t, Point(lch), Point{1.807, 4.532, 214.191}, 3)
}

func TestRGBToXYZConversion(t *testing.T) {
	rd50 := NewRGBTransformer(&SpaceSRGB, &AdaptationBradford, nil, &Scaler8bClamping, 100.0, nil)
	origRgb := RGB{123, 200, 50}

	xyzd50 := rd50.Convert(origRgb)
	testPointAlmostEqual(t, Point(xyzd50), Point{31.334, 46.006, 8.162}, 3)
	rtRgb := rd50.Invert(xyzd50)

	testPointAlmostEqual(t, Point(rtRgb), Point(origRgb), 3)

	// check default whitepoint
	rd50imp := NewRGBTransformer(&SpaceSRGB, &AdaptationBradford, &IlluminantRefD50, &Scaler8bClamping, 100.0, nil)
	xyzd50imp := rd50imp.Convert(origRgb)
	testPointEqual(t, Point(xyzd50imp), Point(xyzd50))

	rd65 := NewRGBTransformer(&SpaceSRGB, &AdaptationBradford, &IlluminantRefD65, &Scaler8bClamping, 100.0, nil)
	xyzd65 := rd65.Convert(origRgb)
	testPointAlmostEqual(t, Point(xyzd65), Point{29.398, 45.748, 10.298}, 3)
}

func TestAdobeRGBToXYZConversion(t *testing.T) {
	rgb := RGB{123, 200, 50}
	rd65 := NewRGBTransformer(&SpaceAdobeRGB, &AdaptationBradford, &IlluminantRefD65, &Scaler8bClamping, 0.0, nil)
	xyzd65 := rd65.Convert(rgb)
	testPointAlmostEqual(t, Point(xyzd65), Point{0.230, 0.430, 0.074}, 3)

	rd50 := NewRGBTransformer(&SpaceAdobeRGB, &AdaptationBradford, nil, &Scaler8bClamping, 0.0, nil)
	xyzd50 := rd50.Convert(rgb)
	testPointAlmostEqual(t, Point(xyzd50), Point{0.247, 0.431, 0.060}, 3)
}

func TestLChToRGBConversion(t *testing.T) {
	lab2xyz := NewLabTransformer(nil)
	rgb2xyz := NewRGBTransformer(&SpaceSRGB, nil, nil, nil, 0.0, nil)
	rgb := rgb2xyz.Invert(lab2xyz.Convert(LCh{0.0, 0.0, 0.0}.Lab()))
	testPointEqual(t, Point(rgb), Point{0.0, 0.0, 0.0})
}

func TestLChuvToLuvConversion(t *testing.T) {
	luv := LChuv{1.807, 2.715, 199.228}.Luv()
	testPointAlmostEqual(t, Point(luv), Point{1.807, -2.564, -0.894}, 3)
}

func TestAdaptation(t *testing.T) {
	xyz := XYZ{0.5, 0.4, 0.1}
	am := AdaptationBradford.Transform(IlluminantRefC.XYZ, IlluminantRefD65.XYZ)
	r := xyz.Transform(&am)
	testPointAlmostEqual(t, Point(r), Point{0.4912, 0.3998, 0.0927}, 4)
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


