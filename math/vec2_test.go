package math

import (
	gomath "math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVec2Add(t *testing.T) {
	v := NewVec2(1, 2).Add(NewVec2(3, 4))
	require.InDelta(t, 4.0, v.X, 1e-9)
	require.InDelta(t, 6.0, v.Y, 1e-9)
}

func TestVec2Sub(t *testing.T) {
	v := NewVec2(5, 7).Sub(NewVec2(2, 3))
	require.InDelta(t, 3.0, v.X, 1e-9)
	require.InDelta(t, 4.0, v.Y, 1e-9)
}

func TestVec2Mul(t *testing.T) {
	v := NewVec2(3, 4).Mul(2)
	require.InDelta(t, 6.0, v.X, 1e-9)
	require.InDelta(t, 8.0, v.Y, 1e-9)
}

func TestVec2Dot(t *testing.T) {
	d := NewVec2(1, 0).Dot(NewVec2(0, 1))
	require.InDelta(t, 0.0, d, 1e-9)

	d = NewVec2(3, 4).Dot(NewVec2(3, 4))
	require.InDelta(t, 25.0, d, 1e-9)
}

func TestVec2Cross(t *testing.T) {
	c := NewVec2(1, 0).Cross(NewVec2(0, 1))
	require.InDelta(t, 1.0, c, 1e-9)
}

func TestVec2Len(t *testing.T) {
	l := NewVec2(3, 4).Len()
	require.InDelta(t, 5.0, l, 1e-9)
}

func TestVec2Normalize(t *testing.T) {
	n := NewVec2(3, 4).Normalize()
	require.InDelta(t, 1.0, n.Len(), 1e-9)

	// Zero vector normalizes to zero
	z := Vec2Zero().Normalize()
	require.Equal(t, Vec2Zero(), z)
}

func TestVec2Distance(t *testing.T) {
	d := NewVec2(0, 0).Distance(NewVec2(3, 4))
	require.InDelta(t, 5.0, d, 1e-9)
}

func TestVec2Lerp(t *testing.T) {
	v := NewVec2(0, 0).Lerp(NewVec2(10, 20), 0.5)
	require.InDelta(t, 5.0, v.X, 1e-9)
	require.InDelta(t, 10.0, v.Y, 1e-9)
}

func TestVec2Rotate(t *testing.T) {
	v := NewVec2(1, 0).Rotate(gomath.Pi / 2)
	require.InDelta(t, 0.0, v.X, 1e-9)
	require.InDelta(t, 1.0, v.Y, 1e-9)
}

func TestVec2Angle(t *testing.T) {
	a := NewVec2(1, 0).Angle()
	require.InDelta(t, 0.0, a, 1e-9)

	a = NewVec2(0, 1).Angle()
	require.InDelta(t, gomath.Pi/2, a, 1e-9)
}

func TestVec2Reflect(t *testing.T) {
	v := NewVec2(1, -1).Reflect(NewVec2(0, 1))
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 1.0, v.Y, 1e-9)
}

func TestVec2Perpendicular(t *testing.T) {
	p := NewVec2(1, 0).Perpendicular()
	require.InDelta(t, 0.0, p.X, 1e-9)
	require.InDelta(t, 1.0, p.Y, 1e-9)
}

func TestVec2Div(t *testing.T) {
	v := NewVec2(6, 8).Div(2)
	require.InDelta(t, 3.0, v.X, 1e-9)
	require.InDelta(t, 4.0, v.Y, 1e-9)
}

func TestVec2MulVec(t *testing.T) {
	v := NewVec2(2, 3).MulVec(NewVec2(4, 5))
	require.InDelta(t, 8.0, v.X, 1e-9)
	require.InDelta(t, 15.0, v.Y, 1e-9)
}

func TestVec2DivVec(t *testing.T) {
	v := NewVec2(10, 15).DivVec(NewVec2(2, 5))
	require.InDelta(t, 5.0, v.X, 1e-9)
	require.InDelta(t, 3.0, v.Y, 1e-9)
}

func TestVec2LenSq(t *testing.T) {
	l := NewVec2(3, 4).LenSq()
	require.InDelta(t, 25.0, l, 1e-9)
}

func TestVec2DistanceSq(t *testing.T) {
	d := NewVec2(0, 0).DistanceSq(NewVec2(3, 4))
	require.InDelta(t, 25.0, d, 1e-9)
}

func TestVec2AngleTo(t *testing.T) {
	a := NewVec2(1, 0).AngleTo(NewVec2(0, 1))
	require.InDelta(t, gomath.Pi/2, a, 1e-9)

	a = NewVec2(0, 1).AngleTo(NewVec2(1, 0))
	require.InDelta(t, -gomath.Pi/2, a, 1e-9)
}

func TestVec2Negate(t *testing.T) {
	v := NewVec2(3, -5).Negate()
	require.InDelta(t, -3.0, v.X, 1e-9)
	require.InDelta(t, 5.0, v.Y, 1e-9)
}

func TestVec2Clamp(t *testing.T) {
	v := NewVec2(-1, 5).Clamp(NewVec2(0, 0), NewVec2(3, 3))
	require.InDelta(t, 0.0, v.X, 1e-9)
	require.InDelta(t, 3.0, v.Y, 1e-9)
}

func TestVec2Floor(t *testing.T) {
	v := NewVec2(1.7, 2.3).Floor()
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 2.0, v.Y, 1e-9)
}

func TestVec2Ceil(t *testing.T) {
	v := NewVec2(1.1, 2.9).Ceil()
	require.InDelta(t, 2.0, v.X, 1e-9)
	require.InDelta(t, 3.0, v.Y, 1e-9)
}

func TestVec2Round(t *testing.T) {
	v := NewVec2(1.4, 2.6).Round()
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 3.0, v.Y, 1e-9)
}

func TestVec2Abs(t *testing.T) {
	v := NewVec2(-3, -4).Abs()
	require.InDelta(t, 3.0, v.X, 1e-9)
	require.InDelta(t, 4.0, v.Y, 1e-9)
}

func TestVec2Min(t *testing.T) {
	v := NewVec2(3, 1).Min(NewVec2(1, 5))
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 1.0, v.Y, 1e-9)
}

func TestVec2Max(t *testing.T) {
	v := NewVec2(3, 1).Max(NewVec2(1, 5))
	require.InDelta(t, 3.0, v.X, 1e-9)
	require.InDelta(t, 5.0, v.Y, 1e-9)
}

func TestVec2String(t *testing.T) {
	s := NewVec2(1, 2).String()
	require.Equal(t, "Vec2(1, 2)", s)
}

func TestVec2DirectionConstructors(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() Vec2
		expected Vec2
	}{
		{"Vec2One", Vec2One, NewVec2(1, 1)},
		{"Vec2Up", Vec2Up, NewVec2(0, -1)},
		{"Vec2Down", Vec2Down, NewVec2(0, 1)},
		{"Vec2Left", Vec2Left, NewVec2(-1, 0)},
		{"Vec2Right", Vec2Right, NewVec2(1, 0)},
		{"Vec2Zero", Vec2Zero, NewVec2(0, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.fn())
		})
	}
}

func TestVec2ApproxEqual(t *testing.T) {
	a := NewVec2(1, 2)
	b := NewVec2(1.0000000001, 2.0000000001)
	require.True(t, a.ApproxEqual(b, 1e-9))
	require.False(t, a.ApproxEqual(NewVec2(2, 2), 1e-9))
}

func BenchmarkVec2Add(b *testing.B) {
	v1 := NewVec2(1, 2)
	v2 := NewVec2(3, 4)
	for b.Loop() {
		_ = v1.Add(v2)
	}
}

func BenchmarkVec2Normalize(b *testing.B) {
	v := NewVec2(3, 4)
	for b.Loop() {
		_ = v.Normalize()
	}
}

func BenchmarkVec2Rotate(b *testing.B) {
	v := NewVec2(1, 0)
	for b.Loop() {
		_ = v.Rotate(0.5)
	}
}
