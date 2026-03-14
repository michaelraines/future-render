package math

import (
	"fmt"
	gomath "math"
)

// Mat4 represents a 4x4 matrix stored in column-major order.
// This is the standard transform type for 3D rendering.
//
// Memory layout (column-major, matching OpenGL/WebGPU conventions):
//
//	| M[0]  M[4]  M[8]  M[12] |
//	| M[1]  M[5]  M[9]  M[13] |
//	| M[2]  M[6]  M[10] M[14] |
//	| M[3]  M[7]  M[11] M[15] |
type Mat4 [16]float64

// Mat4Identity returns the 4x4 identity matrix.
func Mat4Identity() Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Mat4Translate returns a translation matrix.
func Mat4Translate(tx, ty, tz float64) Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		tx, ty, tz, 1,
	}
}

// Mat4Scale returns a scaling matrix.
func Mat4Scale(sx, sy, sz float64) Mat4 {
	return Mat4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
}

// Mat4RotateX returns a rotation matrix around the X axis.
func Mat4RotateX(angle float64) Mat4 {
	sin, cos := gomath.Sincos(angle)
	return Mat4{
		1, 0, 0, 0,
		0, cos, sin, 0,
		0, -sin, cos, 0,
		0, 0, 0, 1,
	}
}

// Mat4RotateY returns a rotation matrix around the Y axis.
func Mat4RotateY(angle float64) Mat4 {
	sin, cos := gomath.Sincos(angle)
	return Mat4{
		cos, 0, -sin, 0,
		0, 1, 0, 0,
		sin, 0, cos, 0,
		0, 0, 0, 1,
	}
}

// Mat4RotateZ returns a rotation matrix around the Z axis.
func Mat4RotateZ(angle float64) Mat4 {
	sin, cos := gomath.Sincos(angle)
	return Mat4{
		cos, sin, 0, 0,
		-sin, cos, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Mat4RotateAxis returns a rotation matrix around an arbitrary axis.
// The axis must be normalized.
func Mat4RotateAxis(axis Vec3, angle float64) Mat4 {
	sin, cos := gomath.Sincos(angle)
	t := 1.0 - cos
	x, y, z := axis.X, axis.Y, axis.Z

	return Mat4{
		t*x*x + cos, t*x*y + sin*z, t*x*z - sin*y, 0,
		t*x*y - sin*z, t*y*y + cos, t*y*z + sin*x, 0,
		t*x*z + sin*y, t*y*z - sin*x, t*z*z + cos, 0,
		0, 0, 0, 1,
	}
}

// Mat4Ortho returns an orthographic projection matrix.
// This maps the box [left,right] × [bottom,top] × [near,far] to [-1,1]³.
func Mat4Ortho(left, right, bottom, top, near, far float64) Mat4 {
	rml := right - left
	tmb := top - bottom
	fmn := far - near
	return Mat4{
		2 / rml, 0, 0, 0,
		0, 2 / tmb, 0, 0,
		0, 0, -2 / fmn, 0,
		-(right + left) / rml, -(top + bottom) / tmb, -(far + near) / fmn, 1,
	}
}

// Mat4Perspective returns a perspective projection matrix.
// fovY is the vertical field of view in radians, aspect is width/height,
// near and far are the clipping plane distances (must be positive).
func Mat4Perspective(fovY, aspect, near, far float64) Mat4 {
	f := 1.0 / gomath.Tan(fovY/2.0)
	nf := near - far
	return Mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / nf, -1,
		0, 0, (2 * far * near) / nf, 0,
	}
}

// Mat4LookAt returns a view matrix looking from eye toward center with the
// given up vector.
func Mat4LookAt(eye, center, up Vec3) Mat4 {
	f := center.Sub(eye).Normalize()
	s := f.Cross(up).Normalize()
	u := s.Cross(f)

	return Mat4{
		s.X, u.X, -f.X, 0,
		s.Y, u.Y, -f.Y, 0,
		s.Z, u.Z, -f.Z, 0,
		-s.Dot(eye), -u.Dot(eye), f.Dot(eye), 1,
	}
}

// At returns the element at row r, column c.
func (m Mat4) At(r, c int) float64 {
	return m[c*4+r]
}

// Set returns a new matrix with the element at row r, column c set to v.
func (m Mat4) Set(r, c int, v float64) Mat4 {
	m[c*4+r] = v
	return m
}

// Col returns column c as a Vec4.
func (m Mat4) Col(c int) Vec4 {
	i := c * 4
	return Vec4{X: m[i], Y: m[i+1], Z: m[i+2], W: m[i+3]}
}

// Row returns row r as a Vec4.
func (m Mat4) Row(r int) Vec4 {
	return Vec4{X: m[r], Y: m[4+r], Z: m[8+r], W: m[12+r]}
}

// Mul returns the product of m and other (m × other).
func (m Mat4) Mul(other Mat4) Mat4 {
	var result Mat4
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			result[c*4+r] = m[r]*other[c*4] +
				m[4+r]*other[c*4+1] +
				m[8+r]*other[c*4+2] +
				m[12+r]*other[c*4+3]
		}
	}
	return result
}

// MulVec4 transforms a Vec4 by this matrix.
func (m Mat4) MulVec4(v Vec4) Vec4 {
	return Vec4{
		X: m[0]*v.X + m[4]*v.Y + m[8]*v.Z + m[12]*v.W,
		Y: m[1]*v.X + m[5]*v.Y + m[9]*v.Z + m[13]*v.W,
		Z: m[2]*v.X + m[6]*v.Y + m[10]*v.Z + m[14]*v.W,
		W: m[3]*v.X + m[7]*v.Y + m[11]*v.Z + m[15]*v.W,
	}
}

// MulVec3Point transforms a 3D point by this matrix (assumes w=1).
func (m Mat4) MulVec3Point(v Vec3) Vec3 {
	w := m[3]*v.X + m[7]*v.Y + m[11]*v.Z + m[15]
	if w == 0 {
		return Vec3{}
	}
	return Vec3{
		X: (m[0]*v.X + m[4]*v.Y + m[8]*v.Z + m[12]) / w,
		Y: (m[1]*v.X + m[5]*v.Y + m[9]*v.Z + m[13]) / w,
		Z: (m[2]*v.X + m[6]*v.Y + m[10]*v.Z + m[14]) / w,
	}
}

// MulVec3Dir transforms a 3D direction by this matrix (assumes w=0).
func (m Mat4) MulVec3Dir(v Vec3) Vec3 {
	return Vec3{
		X: m[0]*v.X + m[4]*v.Y + m[8]*v.Z,
		Y: m[1]*v.X + m[5]*v.Y + m[9]*v.Z,
		Z: m[2]*v.X + m[6]*v.Y + m[10]*v.Z,
	}
}

// Transpose returns the transpose of m.
func (m Mat4) Transpose() Mat4 {
	return Mat4{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

// Determinant returns the determinant of m.
func (m Mat4) Determinant() float64 {
	// Cofactor expansion along the first row
	a := m[5]*(m[10]*m[15]-m[11]*m[14]) - m[9]*(m[6]*m[15]-m[7]*m[14]) + m[13]*(m[6]*m[11]-m[7]*m[10])
	b := m[1]*(m[10]*m[15]-m[11]*m[14]) - m[9]*(m[2]*m[15]-m[3]*m[14]) + m[13]*(m[2]*m[11]-m[3]*m[10])
	c := m[1]*(m[6]*m[15]-m[7]*m[14]) - m[5]*(m[2]*m[15]-m[3]*m[14]) + m[13]*(m[2]*m[7]-m[3]*m[6])
	d := m[1]*(m[6]*m[11]-m[7]*m[10]) - m[5]*(m[2]*m[11]-m[3]*m[10]) + m[9]*(m[2]*m[7]-m[3]*m[6])

	return m[0]*a - m[4]*b + m[8]*c - m[12]*d
}

// Inverse returns the inverse of m and true, or the zero matrix and false
// if m is singular.
func (m Mat4) Inverse() (Mat4, bool) {
	var inv Mat4

	inv[0] = m[5]*m[10]*m[15] - m[5]*m[11]*m[14] - m[9]*m[6]*m[15] + m[9]*m[7]*m[14] + m[13]*m[6]*m[11] - m[13]*m[7]*m[10]
	inv[4] = -m[4]*m[10]*m[15] + m[4]*m[11]*m[14] + m[8]*m[6]*m[15] - m[8]*m[7]*m[14] - m[12]*m[6]*m[11] + m[12]*m[7]*m[10]
	inv[8] = m[4]*m[9]*m[15] - m[4]*m[11]*m[13] - m[8]*m[5]*m[15] + m[8]*m[7]*m[13] + m[12]*m[5]*m[11] - m[12]*m[7]*m[9]
	inv[12] = -m[4]*m[9]*m[14] + m[4]*m[10]*m[13] + m[8]*m[5]*m[14] - m[8]*m[6]*m[13] - m[12]*m[5]*m[10] + m[12]*m[6]*m[9]

	inv[1] = -m[1]*m[10]*m[15] + m[1]*m[11]*m[14] + m[9]*m[2]*m[15] - m[9]*m[3]*m[14] - m[13]*m[2]*m[11] + m[13]*m[3]*m[10]
	inv[5] = m[0]*m[10]*m[15] - m[0]*m[11]*m[14] - m[8]*m[2]*m[15] + m[8]*m[3]*m[14] + m[12]*m[2]*m[11] - m[12]*m[3]*m[10]
	inv[9] = -m[0]*m[9]*m[15] + m[0]*m[11]*m[13] + m[8]*m[1]*m[15] - m[8]*m[3]*m[13] - m[12]*m[1]*m[11] + m[12]*m[3]*m[9]
	inv[13] = m[0]*m[9]*m[14] - m[0]*m[10]*m[13] - m[8]*m[1]*m[14] + m[8]*m[2]*m[13] + m[12]*m[1]*m[10] - m[12]*m[2]*m[9]

	inv[2] = m[1]*m[6]*m[15] - m[1]*m[7]*m[14] - m[5]*m[2]*m[15] + m[5]*m[3]*m[14] + m[13]*m[2]*m[7] - m[13]*m[3]*m[6]
	inv[6] = -m[0]*m[6]*m[15] + m[0]*m[7]*m[14] + m[4]*m[2]*m[15] - m[4]*m[3]*m[14] - m[12]*m[2]*m[7] + m[12]*m[3]*m[6]
	inv[10] = m[0]*m[5]*m[15] - m[0]*m[7]*m[13] - m[4]*m[1]*m[15] + m[4]*m[3]*m[13] + m[12]*m[1]*m[7] - m[12]*m[3]*m[5]
	inv[14] = -m[0]*m[5]*m[14] + m[0]*m[6]*m[13] + m[4]*m[1]*m[14] - m[4]*m[2]*m[13] - m[12]*m[1]*m[6] + m[12]*m[2]*m[5]

	inv[3] = -m[1]*m[6]*m[11] + m[1]*m[7]*m[10] + m[5]*m[2]*m[11] - m[5]*m[3]*m[10] - m[9]*m[2]*m[7] + m[9]*m[3]*m[6]
	inv[7] = m[0]*m[6]*m[11] - m[0]*m[7]*m[10] - m[4]*m[2]*m[11] + m[4]*m[3]*m[10] + m[8]*m[2]*m[7] - m[8]*m[3]*m[6]
	inv[11] = -m[0]*m[5]*m[11] + m[0]*m[7]*m[9] + m[4]*m[1]*m[11] - m[4]*m[3]*m[9] - m[8]*m[1]*m[7] + m[8]*m[3]*m[5]
	inv[15] = m[0]*m[5]*m[10] - m[0]*m[6]*m[9] - m[4]*m[1]*m[10] + m[4]*m[2]*m[9] + m[8]*m[1]*m[6] - m[8]*m[2]*m[5]

	det := m[0]*inv[0] + m[1]*inv[4] + m[2]*inv[8] + m[3]*inv[12]
	if gomath.Abs(det) < 1e-14 {
		return Mat4{}, false
	}

	invDet := 1.0 / det
	for i := range inv {
		inv[i] *= invDet
	}
	return inv, true
}

// Translation returns the translation component of the matrix.
func (m Mat4) Translation() Vec3 {
	return Vec3{X: m[12], Y: m[13], Z: m[14]}
}

// Mat3 returns the upper-left 3x3 submatrix.
func (m Mat4) Mat3() Mat3 {
	return Mat3{
		m[0], m[1], m[2],
		m[4], m[5], m[6],
		m[8], m[9], m[10],
	}
}

// ApproxEqual returns whether m and other are approximately equal within epsilon.
func (m Mat4) ApproxEqual(other Mat4, epsilon float64) bool {
	for i := range m {
		if gomath.Abs(m[i]-other[i]) > epsilon {
			return false
		}
	}
	return true
}

// String returns a human-readable representation of m.
func (m Mat4) String() string {
	return fmt.Sprintf("Mat4[\n  [%g, %g, %g, %g]\n  [%g, %g, %g, %g]\n  [%g, %g, %g, %g]\n  [%g, %g, %g, %g]\n]",
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15])
}
