# go-chromath

go-chromath is a library for color space math including transformations, chromatic adaptation and color difference (ΔE) calculation. It provides implentations of common RGB spaces to CIEXYZ, and CIEXYZ to useful CIE spaces such as L\*a\*b\*, LCH, Luv. Special care has been taken in the implementation of CIEDE2000, a calculation that is commonly defectively implemented.

## Use

Let's convert two colors from sRGB to L\*a\*b with a reference white of D50 and perform a CIEDE2000 color difference calculation. This requires creating a coverter from RGB to XYZ, and from XYZ to Lab. One converter instance can be used on all points.
```go
c1 := RGB{194, 0, 120} // xkcd magenta
c2 := RGB{203, 65, 107} // xkcd hot pink

targetIlluminant := &IlluminantRefD50
```

Create an sRGB Transformer to XYZ with a target illuminant of D50. Since sRGB is defined to have a native illuminant of D65, use the Bradford response transform to adapt.  Use an 8-bit value (0-255) scaler with clamping to scale input values, and set the output scale (for XYZ) to unity. Finally, do not override the native linearization.
```go
rgb2xyz := NewRGBTransformer(&SpaceSRGB, &AdaptationBradford, targetIlluminant, &Scaler8bClamping, 1.0, nil)
```

Create an XYZ ⇔ Lab transformer with an assumed illuminant of D50. Note that this illuminant needs to match the target illuminant for all XYZ transformed points. This will not perform adaptation as the XYZ point objects themselves do not carry their whitepoint information.
```go
lab2xyz := NewLabTransformer(targetIlluminant, 1.0)

// Now converting color points is simple
c1xyz := rgb2xyz.Convert(c1)
c2xyz := rgb2xyz.Convert(c2)
```

Now convert to L\*a\*b. The naming convention is that conversion is "towards" XYZ, and so XYZ ⇒ Any is "inversion"
```go
c1lab := lab2xyz.Invert(c1xyz)
c2lab := lab2xyz.Invert(c2xyz)
```

Perform the ΔE computation, which is provided by a stateless function.
```go
Δe2000 := deltae.CIE2000(c1lab, c2lab, &KLChDefault)
```

Since 8-bit sRGB is such a common model, and the sRGB gamma function is quite expensive in floating point, a faster sRGB scaler compander implementation is provided: `SRGBFastCompander`.
This uses a lookup table for linearization, but is still quite slow for companding, although this could be improved. The compander includes scaling for 8-bit RGB, so should not need a scaler. Use it like this:
```go
// no scaler, and override compander with SRGBFastCompander
rgb2xyz := NewRGBTransformer(&SpaceSRGB, &AdaptationBradford, targetIlluminant, nil, 1.0, SRGBFastCompander.Init(&SpaceSRGB))
```

## TODO
* While there are some basic tests, coverage in transformations needs improvements. The CIEDE2000 is better tested.
* An implementation of CIECAM02 is also forthcoming.
* More documentation is required. Unfortunately all this presupposes a previous familiarlity with the color space and matching fundamentals.
* Continue to refine the packaging
