package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVec3Add(t *testing.T) {
	v := NewVec3(1, 2, 3).Add(NewVec3(4, 5, 6))
	require.InDelta(t, 5.0, v.X, 1e-9)
	require.InDelta(t, 7.0, v.Y, 1e-9)
	require.InDelta(t, 9.0, v.Z, 1e-9)
}

func TestVec3Sub(t *testing.T) {
	v := NewVec3(5, 7, 9).Sub(NewVec3(1, 2, 3))
	require.InDelta(t, 4.0, v.X, 1e-9)
	require.InDelta(t, 5.0, v.Y, 1e-9)
	require.InDelta(t, 6.0, v.Z, 1e-9)
}

func TestVec3Mul(t *testing.T) {
	v := NewVec3(1, 2, 3).Mul(3)
	require.InDelta(t, 3.0, v.X, 1e-9)
	require.InDelta(t, 6.0, v.Y, 1e-9)
	require.InDelta(t, 9.0, v.Z, 1e-9)
}

func TestVec3Div(t *testing.T) {
	v := NewVec3(6, 9, 12).Div(3)
	require.InDelta(t, 2.0, v.X, 1e-9)
	require.InDelta(t, 3.0, v.Y, 1e-9)
	require.InDelta(t, 4.0, v.Z, 1e-9)
}

func TestVec3MulVec(t *testing.T) {
	v := NewVec3(2, 3, 4).MulVec(NewVec3(5, 6, 7))
	require.InDelta(t, 10.0, v.X, 1e-9)
	require.InDelta(t, 18.0, v.Y, 1e-9)
	require.InDelta(t, 28.0, v.Z, 1e-9)
}

func TestVec3Dot(t *testing.T) {
	d := NewVec3(1, 0, 0).Dot(NewVec3(0, 1, 0))
	require.InDelta(t, 0.0, d, 1e-9)

	d = NewVec3(1, 2, 3).Dot(NewVec3(1, 2, 3))
	require.InDelta(t, 14.0, d, 1e-9)
}

func TestVec3Cross(t *testing.T) {
	v := NewVec3(1, 0, 0).Cross(NewVec3(0, 1, 0))
	require.InDelta(t, 0.0, v.X, 1e-9)
	require.InDelta(t, 0.0, v.Y, 1e-9)
	require.InDelta(t, 1.0, v.Z, 1e-9)
}

func TestVec3Len(t *testing.T) {
	l := NewVec3(2, 3, 6).Len()
	require.InDelta(t, 7.0, l, 1e-9)
}

func TestVec3LenSq(t *testing.T) {
	l := NewVec3(2, 3, 6).LenSq()
	require.InDelta(t, 49.0, l, 1e-9)
}

func TestVec3Normalize(t *testing.T) {
	n := NewVec3(3, 4, 0).Normalize()
	require.InDelta(t, 1.0, n.Len(), 1e-9)

	z := Vec3Zero().Normalize()
	require.Equal(t, Vec3Zero(), z)
}

func TestVec3Distance(t *testing.T) {
	d := NewVec3(0, 0, 0).Distance(NewVec3(2, 3, 6))
	require.InDelta(t, 7.0, d, 1e-9)
}

func TestVec3DistanceSq(t *testing.T) {
	d := NewVec3(0, 0, 0).DistanceSq(NewVec3(2, 3, 6))
	require.InDelta(t, 49.0, d, 1e-9)
}

func TestVec3Lerp(t *testing.T) {
	v := NewVec3(0, 0, 0).Lerp(NewVec3(10, 20, 30), 0.5)
	require.InDelta(t, 5.0, v.X, 1e-9)
	require.InDelta(t, 10.0, v.Y, 1e-9)
	require.InDelta(t, 15.0, v.Z, 1e-9)
}

func TestVec3Negate(t *testing.T) {
	v := NewVec3(1, -2, 3).Negate()
	require.InDelta(t, -1.0, v.X, 1e-9)
	require.InDelta(t, 2.0, v.Y, 1e-9)
	require.InDelta(t, -3.0, v.Z, 1e-9)
}

func TestVec3Reflect(t *testing.T) {
	v := NewVec3(1, -1, 0).Reflect(NewVec3(0, 1, 0))
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 1.0, v.Y, 1e-9)
	require.InDelta(t, 0.0, v.Z, 1e-9)
}

func TestVec3Clamp(t *testing.T) {
	v := NewVec3(-1, 5, 0.5).Clamp(NewVec3(0, 0, 0), NewVec3(3, 3, 3))
	require.InDelta(t, 0.0, v.X, 1e-9)
	require.InDelta(t, 3.0, v.Y, 1e-9)
	require.InDelta(t, 0.5, v.Z, 1e-9)
}

func TestVec3Abs(t *testing.T) {
	v := NewVec3(-1, -2, -3).Abs()
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 2.0, v.Y, 1e-9)
	require.InDelta(t, 3.0, v.Z, 1e-9)
}

func TestVec3Min(t *testing.T) {
	v := NewVec3(3, 1, 5).Min(NewVec3(1, 5, 2))
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 1.0, v.Y, 1e-9)
	require.InDelta(t, 2.0, v.Z, 1e-9)
}

func TestVec3Max(t *testing.T) {
	v := NewVec3(3, 1, 5).Max(NewVec3(1, 5, 2))
	require.InDelta(t, 3.0, v.X, 1e-9)
	require.InDelta(t, 5.0, v.Y, 1e-9)
	require.InDelta(t, 5.0, v.Z, 1e-9)
}

func TestVec3String(t *testing.T) {
	s := NewVec3(1, 2, 3).String()
	require.Equal(t, "Vec3(1, 2, 3)", s)
}

func TestVec3Vec2(t *testing.T) {
	v := NewVec3(1, 2, 3).Vec2()
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 2.0, v.Y, 1e-9)
}

func TestVec3ApproxEqual(t *testing.T) {
	a := NewVec3(1, 2, 3)
	b := NewVec3(1.0000000001, 2.0000000001, 3.0000000001)
	require.True(t, a.ApproxEqual(b, 1e-9))
	require.False(t, a.ApproxEqual(NewVec3(2, 2, 3), 1e-9))
}

func TestVec3Constructors(t *testing.T) {
	require.Equal(t, NewVec3(0, 0, 0), Vec3Zero())
	require.Equal(t, NewVec3(1, 1, 1), Vec3One())
	require.Equal(t, NewVec3(1, 0, 0), Vec3UnitX())
	require.Equal(t, NewVec3(0, 1, 0), Vec3UnitY())
	require.Equal(t, NewVec3(0, 0, 1), Vec3UnitZ())
}
