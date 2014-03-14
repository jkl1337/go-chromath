package chromath

const Gamma22 = 563/256.0

type RGBSpace struct {
	Name        string
	Description string
	XyYPrimary XyYPrimary
	IlluminantRef *IlluminantRef
	Gamma Gamma
	Compander Compander
}

var SpaceSRGB = RGBSpace{
	"sRGB", "sRGB",
	XyYPrimary{
		0.6400, 0.3300,
		0.3000, 0.6000,
		0.1500, 0.0600,
	},
	&IlluminantRefD65,
	Gamma22,
	&SRGBCompander,
}

var SpaceAdobeRGB = RGBSpace{
	"Adobe RGB", "Adobe RGB (1998)",
	XyYPrimary{
		0.6400, 0.3300,
		0.2100, 0.7100,
		0.1500, 0.0600,
	},
	&IlluminantRefD65,
	Gamma22,
	&GammaCompander,
}

var SpaceAppleRGB = RGBSpace{
	"Apple RGB", "Apple RGB",
	XyYPrimary{
		0.6250, 0.3400,
		0.2800,	0.5950,
		0.1550,	0.0700,
	},
	&IlluminantRefD65,
	1.8,
	&GammaCompander,
}

var SpaceBestRGB = RGBSpace{
	"Best RGB", "Best RGB",
	XyYPrimary{
		0.7347, 0.2653,
		0.2150,	0.7750,
		0.1300,	0.0350,
	},
	&IlluminantRefD50,
	Gamma22,
	&GammaCompander,
}

var SpaceBetaRGB = RGBSpace{
	"Beta RGB", "Beta RGB",
	XyYPrimary{
		0.6888, 0.3112,
		0.1986,	0.7551,
		0.1265,	0.0352,
	},
	&IlluminantRefD50,
	Gamma22,
	&GammaCompander,
}

var SpaceBruceRGB = RGBSpace{
	"Bruce RGB", "Bruce RGB",
	XyYPrimary{
		0.6400, 0.3300,
		0.2800,	0.6500,
		0.1500,	0.0600,
	},
	&IlluminantRefD65,
	Gamma22,
	&GammaCompander,
}

var SpaceCIERGB = RGBSpace{
	"CIE RGB", "CIE RGB",
	XyYPrimary{
		0.7350,	0.2650,
		0.2740,	0.7170,
		0.1670,	0.0090,
	},
	&IlluminantRefE,
	Gamma22,
	&GammaCompander,
}

var SpaceColorMatchRGB = RGBSpace{
	"ColorMatch RGB", "ColorMatch RGB",
	XyYPrimary{
		0.6300,	0.3400,
		0.2950,	0.6050,
		0.1500, 0.0750,
	},
	&IlluminantRefD50,
	1.8,
	&GammaCompander,
}

var SpaceECIRGB = RGBSpace{
	"ECI RGB", "ECI RGB v2",
	XyYPrimary{
		0.6700,	0.3300,
		0.2100,	0.7100,
		0.1400,	0.0800,
	},
	&IlluminantRefD50,
	0.0,
	&LstarCompander,
}

var SpaceNTSCRGB = RGBSpace{
	"NTSC RGB", "NTSC (1953) RGB",
	XyYPrimary{
		0.6700,	0.3300,
		0.2100,	0.7100,
		0.1400,	0.0800,
	},
	&IlluminantRefC,
	2.2,
	&GammaCompander,
}

var SpacePALSECAMRGB = RGBSpace{
	"PAL/SECAM RGB", "PAL/SECAM RGB",
	XyYPrimary{
		0.6400,	0.3300,
		0.2900,	0.6000,
		0.1500,	0.0600,
	},
	&IlluminantRefD65,
	2.2,
	&GammaCompander,
}

var SpaceProPhotoRGB = RGBSpace{
	"ProPhoto RGB", "ProPhoto RGB",
	XyYPrimary{
		0.7347,	0.2653,
		0.1596,	0.8404,
		0.0366,	0.0001,
	},
	&IlluminantRefD50,
	1.8,
	&GammaCompander,
}

var SpaceSMPTECRGB = RGBSpace{
	"SMPTE-C RGB", "SMPTE-C (NTSC 1987) RGB",
	XyYPrimary{
		0.6300,	0.3400,
		0.3100,	0.5950,
		0.1550,	0.0700,
	},
	&IlluminantRefD65,
	2.2,
	&GammaCompander,
}

var SpaceUHDTVRGB = RGBSpace{
	"UHDTV RGB", "UHDTV (ITU-R BT.2020) RGB",
	XyYPrimary{
		0.708, 0.292,
		0.170, 0.797,
		0.131, 0.046,
	},
	&IlluminantRefD65,
	2.4,
	&BT2020Compander,
}

var SpaceWideGamutRGB = RGBSpace{
	"Wide Gamut RGB", "Adobe Wide Gamut RGB",
	XyYPrimary{
		0.7350,	0.2650,
		0.1150,	0.8260,
		0.1570,	0.0180,
	},
	&IlluminantRefD50,
	1.8,
	&GammaCompander,
}

var AdaptationBradford = Adaptation{
	0.8951, -0.7502, 0.03890,
	0.2664, 1.7135, -0.06850,
	-0.1614000, 0.03670, 1.02960,
}

var AdaptationXYZScaling = Adaptation{1, 0, 0, 0, 1, 0, 0, 0, 1}

var AdaptationVonKries = Adaptation{
	0.400240, -0.22630, 0.000,
	0.707600, 1.1653200, 0.000,
	-0.0808100, 0.045700, 0.918220,
}

type Observer int

const (
	CIE2 Observer = iota
	CIE10
)

type IlluminantRef struct {
	XYZ XYZ
	Observer Observer
	Standard *IlluminantStd
}

type IlluminantStd struct {
	Name        string
	Description string
	CCT         int
}

var IlluminantA = IlluminantStd{"A", "Incandescent / Tungsten", 2856}
var IlluminantB = IlluminantStd{"B", "Direct sunlight at noon (obsolete)", 4874}
var IlluminantC = IlluminantStd{"C", "Average / North sky Daylight (obsolete)", 6774}
var IlluminantD50 = IlluminantStd{"D50", "Horizon Light, ICC profile PCS", 5003}
var IlluminantD55 = IlluminantStd{"D55", "Mid-morning / Mid-afternoon Daylight", 5503}
var IlluminantD65 = IlluminantStd{"D65", "Noon Daylight / Television / sRGB", 6504}
var IlluminantD75 = IlluminantStd{"D75", "North sky Daylight", 7504}
var IlluminantE = IlluminantStd{"E", "Equal energy", 5454}
var IlluminantF2 = IlluminantStd{"F2", "Cool White Fluorescent", 4230}
var IlluminantF7 = IlluminantStd{"F7", "D65 simulator / Daylight Simulator", 6500}
var IlluminantF11 = IlluminantStd{"F11", "Philips TL84, Ultralum 40", 4000}

// CIE 1931 Standard (2°) observer illuminant tristimulus values
var IlluminantRefA = IlluminantRef{XYZ{1.09850, 1.00000, 0.35585}, CIE2, &IlluminantA}
var IlluminantRefB = IlluminantRef{XYZ{0.99072, 1.00000, 0.85223}, CIE2, &IlluminantB}
var IlluminantRefC = IlluminantRef{XYZ{0.98074, 1.00000, 1.18232}, CIE2, &IlluminantC}
var IlluminantRefD50 = IlluminantRef{XYZ{0.96422, 1.00000, 0.82521}, CIE2, &IlluminantD50}
var IlluminantRefD55 = IlluminantRef{XYZ{0.95682, 1.00000, 0.92149}, CIE2, &IlluminantD55}
var IlluminantRefD65 = IlluminantRef{XYZ{0.95047, 1.00000, 1.08883}, CIE2, &IlluminantD65}
var IlluminantRefD75 = IlluminantRef{XYZ{0.94972, 1.00000, 1.22638}, CIE2, &IlluminantD75}
var IlluminantRefE = IlluminantRef{XYZ{1.00000, 1.00000, 1.00000}, CIE2, &IlluminantE}
var IlluminantRefF2 = IlluminantRef{XYZ{0.99186, 1.00000, 0.67393}, CIE2, &IlluminantF2}
var IlluminantRefF7 = IlluminantRef{XYZ{0.95041, 1.00000, 1.08747}, CIE2, &IlluminantF7}
var IlluminantRefF11 = IlluminantRef{XYZ{1.00962, 1.00000, 0.64350}, CIE2, &IlluminantF11}

// Supplementary (10°) observer illuminant tristimulus values
// source Berns 2000, cross-checked with python-colormath
var IlluminantRefSuppA = IlluminantRef{XYZ{1.1114, 1.0000, 0.3520}, CIE10, &IlluminantA}
var IlluminantRefSuppC = IlluminantRef{XYZ{0.9728, 1.0000, 1.11614}, CIE10, &IlluminantC}
var IlluminantRefSuppD50 = IlluminantRef{XYZ{0.9672, 1.0000, 0.8143}, CIE10, &IlluminantD50}
var IlluminantRefSuppD55 = IlluminantRef{XYZ{0.9580, 1.0000, 0.9093}, CIE10, &IlluminantD55}
var IlluminantRefSuppD65 = IlluminantRef{XYZ{0.9481, 1.0000, 1.0730}, CIE10, &IlluminantD65}
var IlluminantRefSuppD75 = IlluminantRef{XYZ{0.94416, 1.0000, 1.2064}, CIE10, &IlluminantD75}
var IlluminantRefSuppF2 = IlluminantRef{XYZ{1.0328, 1.0000, 0.6902}, CIE10, &IlluminantF2}
