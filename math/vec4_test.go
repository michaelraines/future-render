package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVec4Add(t *testing.T) {
	v := NewVec4(1, 2, 3, 4).Add(NewVec4(5, 6, 7, 8))
	require.InDelta(t, 6.0, v.X, 1e-9)
	require.InDelta(t, 8.0, v.Y, 1e-9)
	require.InDelta(t, 10.0, v.Z, 1e-9)
	require.InDelta(t, 12.0, v.W, 1e-9)
}

func TestVec4Sub(t *testing.T) {
	v := NewVec4(5, 7, 9, 11).Sub(NewVec4(1, 2, 3, 4))
	require.InDelta(t, 4.0, v.X, 1e-9)
	require.InDelta(t, 5.0, v.Y, 1e-9)
	require.InDelta(t, 6.0, v.Z, 1e-9)
	require.InDelta(t, 7.0, v.W, 1e-9)
}

func TestVec4Mul(t *testing.T) {
	v := NewVec4(1, 2, 3, 4).Mul(2)
	require.InDelta(t, 2.0, v.X, 1e-9)
	require.InDelta(t, 4.0, v.Y, 1e-9)
	require.InDelta(t, 6.0, v.Z, 1e-9)
	require.InDelta(t, 8.0, v.W, 1e-9)
}

func TestVec4Dot(t *testing.T) {
	d := NewVec4(1, 2, 3, 4).Dot(NewVec4(1, 2, 3, 4))
	require.InDelta(t, 30.0, d, 1e-9)
}

func TestVec4Len(t *testing.T) {
	l := NewVec4(1, 0, 0, 0).Len()
	require.InDelta(t, 1.0, l, 1e-9)
}

func TestVec4LenSq(t *testing.T) {
	l := NewVec4(1, 2, 3, 4).LenSq()
	require.InDelta(t, 30.0, l, 1e-9)
}

func TestVec4Normalize(t *testing.T) {
	n := NewVec4(3, 4, 0, 0).Normalize()
	require.InDelta(t, 1.0, n.Len(), 1e-9)

	z := Vec4Zero().Normalize()
	require.Equal(t, Vec4Zero(), z)
}

func TestVec4Lerp(t *testing.T) {
	v := NewVec4(0, 0, 0, 0).Lerp(NewVec4(10, 20, 30, 40), 0.5)
	require.InDelta(t, 5.0, v.X, 1e-9)
	require.InDelta(t, 10.0, v.Y, 1e-9)
	require.InDelta(t, 15.0, v.Z, 1e-9)
	require.InDelta(t, 20.0, v.W, 1e-9)
}

func TestVec4String(t *testing.T) {
	s := NewVec4(1, 2, 3, 4).String()
	require.Equal(t, "Vec4(1, 2, 3, 4)", s)
}

func TestVec4Vec3(t *testing.T) {
	v := NewVec4(1, 2, 3, 4).Vec3()
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 2.0, v.Y, 1e-9)
	require.InDelta(t, 3.0, v.Z, 1e-9)
}

func TestVec4Vec2(t *testing.T) {
	v := NewVec4(1, 2, 3, 4).Vec2()
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 2.0, v.Y, 1e-9)
}

func TestVec4PerspectiveDivide(t *testing.T) {
	v := NewVec4(2, 4, 6, 2).PerspectiveDivide()
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 2.0, v.Y, 1e-9)
	require.InDelta(t, 3.0, v.Z, 1e-9)

	// W=0 returns zero
	z := NewVec4(1, 2, 3, 0).PerspectiveDivide()
	require.Equal(t, Vec3Zero(), z)
}

func TestVec4ApproxEqual(t *testing.T) {
	a := NewVec4(1, 2, 3, 4)
	b := NewVec4(1.0000000001, 2.0000000001, 3.0000000001, 4.0000000001)
	require.True(t, a.ApproxEqual(b, 1e-9))
	require.False(t, a.ApproxEqual(NewVec4(2, 2, 3, 4), 1e-9))
}

func TestVec4Zero(t *testing.T) {
	require.Equal(t, NewVec4(0, 0, 0, 0), Vec4Zero())
}
