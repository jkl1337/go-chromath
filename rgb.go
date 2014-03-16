package chromath

// RGBTransformer allows tranforms from a user-defined RGB color space to XYZ with optional scaling and adaptation
type RGBTransformer struct {
	ws                      *RGBSpace
	compander               Compander
	transform, transformInv Matrix
	inScaler                Scaler
	outScale                float64
	spaceTm, adaptTm        Matrix
}

// NewRGBTransformer initializes a transform operator from an RGB space to XYZ
// ws specifies the parameters (primaries, working illuminant, and linearization function) of the working RGB space
// adaptation specifies an optional adaptation method to transform into the target illuminant.
// targetIluminant specifies the target illuminant of the XYZ conversion. It is ignored if adaptation is nil.
// The default is CIE D50 (the ICC Profile Connection illuminant)
// scaler is an optional scaling function
// outscale provides a scaling factor into XYZ space, typical values are 1 and 100
// compander is an optional specification for the gamma function to linearize RGB points. This argument overrides
// any compander specified by ws.
func NewRGBTransformer(ws *RGBSpace, adaptation *Adaptation, targetIlluminant *IlluminantRef,
	scaler Scaler, outScale float64,
	compander Compander) *RGBTransformer {
	t := &RGBTransformer{}
	t.ws = ws

	if compander != nil {
		t.compander = compander
	} else if ws.Compander != nil {
		t.compander = ws.Compander.Init(ws)
	}

	t.inScaler = scaler

	t.spaceTm = ws.XyYPrimary.RGBTransform(ws.IlluminantRef.XYZ)

	x := t.spaceTm

	if adaptation != nil {
		if targetIlluminant == nil {
			targetIlluminant = &IlluminantRefD50
		}
		if targetIlluminant.XYZ != ws.IlluminantRef.XYZ {
			t.adaptTm = adaptation.Transform(ws.IlluminantRef.XYZ, targetIlluminant.XYZ)
			x = t.adaptTm.Mul3(x)
		}
	}

	if outScale != 0 {
		t.outScale = outScale
		x = x.Mul3(Matrix{outScale, 0, 0, 0, outScale, 0, 0, 0, outScale})
	}
	t.transform = x
	t.transformInv = x.Inv()

	return t
}

// Linearize will apply the current scaling and linearization function to the given RGB point
func (t *RGBTransformer) Linearize(rgb RGB) RGB {
	p := Point(rgb)
	if scaler := t.inScaler; scaler != nil {
		p = scaler.Scale(p)
	}
	if compander := t.compander; compander != nil {
		p = compander.Linearize(p)
	}
	return RGB(p)
}

// Convert an RGB colorpoint to XYZ based on the parameters of the receiver
func (t *RGBTransformer) Convert(rgb RGB) XYZ {
	p := Point(rgb)
	if scaler := t.inScaler; scaler != nil {
		p = scaler.Scale(p)
	}
	if compander := t.compander; compander != nil {
		p = compander.Linearize(p)
	}
	p = t.transform.Mul3x1(p)

	return XYZ(p)
}

// Invert converts an XYZ point to RGB based on the parameters of the receiver
func (t *RGBTransformer) Invert(xyz XYZ) RGB {
	p := Point(xyz)
	p = t.transformInv.Mul3x1(p)

	if compander := t.compander; compander != nil {
		p = compander.Compand(p)
	}

	if scaler := t.inScaler; scaler != nil {
		p = scaler.ScaleInv(p)
	}

	return RGB(p)
}

// Return the computed adaptation matrix for the transformer
func (t *RGBTransformer) AdaptTM() Matrix {
	return t.adaptTm
}

// Return the computed linear RGB⇒XYZ matrix for the transformer
func (t *RGBTransformer) SpaceTM() Matrix {
	return t.spaceTm
}

// Return the computed complete transformation matrix for the transformer
// This matrix includes adaptation, space transform, and output scaling
func (t *RGBTransformer) TM() Matrix {
	return t.transform
}
