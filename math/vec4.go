package math

import (
	"fmt"
	gomath "math"
)

// Vec4 represents a 4D vector, commonly used for homogeneous coordinates
// and color values (RGBA).
type Vec4 struct {
	X, Y, Z, W float64
}

// NewVec4 creates a Vec4 from x, y, z, w components.
func NewVec4(x, y, z, w float64) Vec4 {
	return Vec4{X: x, Y: y, Z: z, W: w}
}

// Vec4Zero returns the zero vector.
func Vec4Zero() Vec4 { return Vec4{} }

// Vec3 returns the XYZ components as a Vec3.
func (v Vec4) Vec3() Vec3 {
	return Vec3{X: v.X, Y: v.Y, Z: v.Z}
}

// Vec2 returns the XY components as a Vec2.
func (v Vec4) Vec2() Vec2 {
	return Vec2{X: v.X, Y: v.Y}
}

// PerspectiveDivide returns the Vec3 result of dividing XYZ by W.
// This is used to convert from homogeneous to Cartesian coordinates.
func (v Vec4) PerspectiveDivide() Vec3 {
	if v.W == 0 {
		return Vec3{}
	}
	return Vec3{X: v.X / v.W, Y: v.Y / v.W, Z: v.Z / v.W}
}

// Add returns the sum of v and other.
func (v Vec4) Add(other Vec4) Vec4 {
	return Vec4{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z, W: v.W + other.W}
}

// Sub returns the difference of v and other.
func (v Vec4) Sub(other Vec4) Vec4 {
	return Vec4{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z, W: v.W - other.W}
}

// Mul returns v scaled by s.
func (v Vec4) Mul(s float64) Vec4 {
	return Vec4{X: v.X * s, Y: v.Y * s, Z: v.Z * s, W: v.W * s}
}

// Dot returns the dot product of v and other.
func (v Vec4) Dot(other Vec4) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z + v.W*other.W
}

// Len returns the length of v.
func (v Vec4) Len() float64 {
	return gomath.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W)
}

// LenSq returns the squared length of v.
func (v Vec4) LenSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
}

// Normalize returns a unit vector in the direction of v.
func (v Vec4) Normalize() Vec4 {
	l := v.Len()
	if l == 0 {
		return Vec4{}
	}
	return Vec4{X: v.X / l, Y: v.Y / l, Z: v.Z / l, W: v.W / l}
}

// Lerp returns the linear interpolation between v and other by t.
func (v Vec4) Lerp(other Vec4, t float64) Vec4 {
	return Vec4{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
		Z: v.Z + (other.Z-v.Z)*t,
		W: v.W + (other.W-v.W)*t,
	}
}

// ApproxEqual returns whether v and other are approximately equal within epsilon.
func (v Vec4) ApproxEqual(other Vec4, epsilon float64) bool {
	return gomath.Abs(v.X-other.X) <= epsilon &&
		gomath.Abs(v.Y-other.Y) <= epsilon &&
		gomath.Abs(v.Z-other.Z) <= epsilon &&
		gomath.Abs(v.W-other.W) <= epsilon
}

// String returns a human-readable representation of v.
func (v Vec4) String() string {
	return fmt.Sprintf("Vec4(%g, %g, %g, %g)", v.X, v.Y, v.Z, v.W)
}
