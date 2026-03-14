package math

import (
	"fmt"
	gomath "math"
)

// Mat3 represents a 3x3 matrix stored in column-major order.
// This is the primary transform type for 2D rendering (affine transforms).
//
// Memory layout (column-major):
//
//	| M[0] M[3] M[6] |
//	| M[1] M[4] M[7] |
//	| M[2] M[5] M[8] |
type Mat3 [9]float64

// Mat3Identity returns the 3x3 identity matrix.
func Mat3Identity() Mat3 {
	return Mat3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}

// Mat3FromRows creates a Mat3 from row values (converted to column-major storage).
// This is convenient for human-readable initialization:
//
//	Mat3FromRows(
//	    a, b, c,   // row 0
//	    d, e, f,   // row 1
//	    g, h, i,   // row 2
//	)
func Mat3FromRows(
	m00, m01, m02,
	m10, m11, m12,
	m20, m21, m22 float64,
) Mat3 {
	return Mat3{
		m00, m10, m20,
		m01, m11, m21,
		m02, m12, m22,
	}
}

// Mat3Translate returns a 2D translation matrix.
func Mat3Translate(tx, ty float64) Mat3 {
	return Mat3{
		1, 0, 0,
		0, 1, 0,
		tx, ty, 1,
	}
}

// Mat3Scale returns a 2D scaling matrix.
func Mat3Scale(sx, sy float64) Mat3 {
	return Mat3{
		sx, 0, 0,
		0, sy, 0,
		0, 0, 1,
	}
}

// Mat3Rotate returns a 2D rotation matrix for the given angle in radians.
func Mat3Rotate(angle float64) Mat3 {
	sin, cos := gomath.Sincos(angle)
	return Mat3{
		cos, sin, 0,
		-sin, cos, 0,
		0, 0, 1,
	}
}

// Mat3Shear returns a 2D shear matrix.
func Mat3Shear(sx, sy float64) Mat3 {
	return Mat3{
		1, sy, 0,
		sx, 1, 0,
		0, 0, 1,
	}
}

// At returns the element at row r, column c.
func (m Mat3) At(r, c int) float64 {
	return m[c*3+r]
}

// Set returns a new matrix with the element at row r, column c set to v.
func (m Mat3) Set(r, c int, v float64) Mat3 {
	m[c*3+r] = v
	return m
}

// Col returns column c as a Vec3.
func (m Mat3) Col(c int) Vec3 {
	i := c * 3
	return Vec3{X: m[i], Y: m[i+1], Z: m[i+2]}
}

// Row returns row r as a Vec3.
func (m Mat3) Row(r int) Vec3 {
	return Vec3{X: m[r], Y: m[3+r], Z: m[6+r]}
}

// Mul returns the product of m and other (m × other).
func (m Mat3) Mul(other Mat3) Mat3 {
	var result Mat3
	for c := 0; c < 3; c++ {
		for r := 0; r < 3; r++ {
			result[c*3+r] = m[r]*other[c*3] +
				m[3+r]*other[c*3+1] +
				m[6+r]*other[c*3+2]
		}
	}
	return result
}

// MulVec2 transforms a 2D point by this matrix (assumes w=1).
func (m Mat3) MulVec2(v Vec2) Vec2 {
	return Vec2{
		X: m[0]*v.X + m[3]*v.Y + m[6],
		Y: m[1]*v.X + m[4]*v.Y + m[7],
	}
}

// MulVec2Dir transforms a 2D direction by this matrix (assumes w=0,
// translation is not applied).
func (m Mat3) MulVec2Dir(v Vec2) Vec2 {
	return Vec2{
		X: m[0]*v.X + m[3]*v.Y,
		Y: m[1]*v.X + m[4]*v.Y,
	}
}

// MulVec3 transforms a Vec3 by this matrix.
func (m Mat3) MulVec3(v Vec3) Vec3 {
	return Vec3{
		X: m[0]*v.X + m[3]*v.Y + m[6]*v.Z,
		Y: m[1]*v.X + m[4]*v.Y + m[7]*v.Z,
		Z: m[2]*v.X + m[5]*v.Y + m[8]*v.Z,
	}
}

// Transpose returns the transpose of m.
func (m Mat3) Transpose() Mat3 {
	return Mat3{
		m[0], m[3], m[6],
		m[1], m[4], m[7],
		m[2], m[5], m[8],
	}
}

// Determinant returns the determinant of m.
func (m Mat3) Determinant() float64 {
	return m[0]*(m[4]*m[8]-m[7]*m[5]) -
		m[3]*(m[1]*m[8]-m[7]*m[2]) +
		m[6]*(m[1]*m[5]-m[4]*m[2])
}

// Inverse returns the inverse of m and true, or the zero matrix and false
// if m is singular.
func (m Mat3) Inverse() (Mat3, bool) {
	det := m.Determinant()
	if gomath.Abs(det) < 1e-14 {
		return Mat3{}, false
	}

	invDet := 1.0 / det

	return Mat3{
		(m[4]*m[8] - m[5]*m[7]) * invDet,
		(m[2]*m[7] - m[1]*m[8]) * invDet,
		(m[1]*m[5] - m[2]*m[4]) * invDet,
		(m[5]*m[6] - m[3]*m[8]) * invDet,
		(m[0]*m[8] - m[2]*m[6]) * invDet,
		(m[2]*m[3] - m[0]*m[5]) * invDet,
		(m[3]*m[7] - m[4]*m[6]) * invDet,
		(m[1]*m[6] - m[0]*m[7]) * invDet,
		(m[0]*m[4] - m[1]*m[3]) * invDet,
	}, true
}

// ApproxEqual returns whether m and other are approximately equal within epsilon.
func (m Mat3) ApproxEqual(other Mat3, epsilon float64) bool {
	for i := range m {
		if gomath.Abs(m[i]-other[i]) > epsilon {
			return false
		}
	}
	return true
}

// String returns a human-readable representation of m.
func (m Mat3) String() string {
	return fmt.Sprintf("Mat3[\n  [%g, %g, %g]\n  [%g, %g, %g]\n  [%g, %g, %g]\n]",
		m[0], m[3], m[6],
		m[1], m[4], m[7],
		m[2], m[5], m[8])
}
