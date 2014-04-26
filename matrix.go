package chromath

// Matrix represents a 3x3 matrix in floating point suitable for affine transforms
type Matrix [9]float64

// Mul multiplies the receiver with a scalar
func (a Matrix) Mul(b float64) Matrix {
	return Matrix{a[0]*b, a[1]*b, a[2]*b, a[3]*b, a[4]*b, a[5]*b, a[6]*b, a[7]*b, a[8]*b}
}

// Mul3x1 returns the product of the receiver with the passed in point, which is treated as a column vector.
func (a Matrix) Mul3x1(b Point) Point {
	return Point{
		a[0]*b[0] + a[3]*b[1] + a[6]*b[2],
		a[1]*b[0] + a[4]*b[1] + a[7]*b[2],
		a[2]*b[0] + a[5]*b[1] + a[8]*b[2],
	}
}

// Mul3 returns the matrix product of the receiver and the parameter a.Mul(b) returns b*a.
func (a Matrix) Mul3(b Matrix) Matrix {
	return Matrix{
		a[0]*b[0] + a[3]*b[1] + a[6]*b[2], a[1]*b[0] + a[4]*b[1] + a[7]*b[2], a[2]*b[0] + a[5]*b[1] + a[8]*b[2],
		a[0]*b[3] + a[3]*b[4] + a[6]*b[5], a[1]*b[3] + a[4]*b[4] + a[7]*b[5], a[2]*b[3] + a[5]*b[4] + a[8]*b[5],
		a[0]*b[6] + a[3]*b[7] + a[6]*b[8], a[1]*b[6] + a[4]*b[7] + a[7]*b[8], a[2]*b[6] + a[5]*b[7] + a[8]*b[8],
	}
}

// Transpose computes and returns the matrix transpose of the receiver
func (a Matrix) Transpose() Matrix {
	return Matrix{a[0], a[3], a[6], a[1], a[4], a[7], a[2], a[5], a[8]}
}

// Det computes the determinant of the receiver.
func (a Matrix) Det() float64 {
	return a[0]*a[4]*a[8] + a[3]*a[7]*a[2] + a[6]*a[1]*a[5] - a[6]*a[4]*a[2] - a[3]*a[1]*a[8] - a[0]*a[7]*a[5]
}

// Inv computes the inverse of the receiver. If the inverse is undefined, a runtime panic will
// occur.
func (a Matrix) Inv() Matrix {
	det := a.Det()
	return Matrix{
		a[4]*a[8] - a[5]*a[7], a[2]*a[7] - a[1]*a[8], a[1]*a[5] - a[2]*a[4],
		a[5]*a[6] - a[3]*a[8], a[0]*a[8] - a[2]*a[6], a[2]*a[3] - a[0]*a[5],
		a[3]*a[7] - a[4]*a[6], a[1]*a[6] - a[0]*a[7], a[0]*a[4] - a[1]*a[3],
	}.Mul(1 / det)
}
