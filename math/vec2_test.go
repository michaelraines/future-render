package math

import (
	gomath "math"
	"testing"
)

const testEpsilon = 1e-10

func TestVec2Add(t *testing.T) {
	v := NewVec2(1, 2).Add(NewVec2(3, 4))
	if !v.ApproxEqual(NewVec2(4, 6), testEpsilon) {
		t.Errorf("expected (4, 6), got %v", v)
	}
}

func TestVec2Sub(t *testing.T) {
	v := NewVec2(5, 7).Sub(NewVec2(2, 3))
	if !v.ApproxEqual(NewVec2(3, 4), testEpsilon) {
		t.Errorf("expected (3, 4), got %v", v)
	}
}

func TestVec2Mul(t *testing.T) {
	v := NewVec2(3, 4).Mul(2)
	if !v.ApproxEqual(NewVec2(6, 8), testEpsilon) {
		t.Errorf("expected (6, 8), got %v", v)
	}
}

func TestVec2Dot(t *testing.T) {
	d := NewVec2(1, 0).Dot(NewVec2(0, 1))
	if d != 0 {
		t.Errorf("expected 0, got %g", d)
	}
	d = NewVec2(3, 4).Dot(NewVec2(3, 4))
	if !ApproxEqual(d, 25, testEpsilon) {
		t.Errorf("expected 25, got %g", d)
	}
}

func TestVec2Cross(t *testing.T) {
	c := NewVec2(1, 0).Cross(NewVec2(0, 1))
	if !ApproxEqual(c, 1, testEpsilon) {
		t.Errorf("expected 1, got %g", c)
	}
}

func TestVec2Len(t *testing.T) {
	l := NewVec2(3, 4).Len()
	if !ApproxEqual(l, 5, testEpsilon) {
		t.Errorf("expected 5, got %g", l)
	}
}

func TestVec2Normalize(t *testing.T) {
	n := NewVec2(3, 4).Normalize()
	if !ApproxEqual(n.Len(), 1, testEpsilon) {
		t.Errorf("expected length 1, got %g", n.Len())
	}
	// Zero vector normalizes to zero
	z := Vec2Zero().Normalize()
	if z != Vec2Zero() {
		t.Errorf("expected zero vector, got %v", z)
	}
}

func TestVec2Distance(t *testing.T) {
	d := NewVec2(0, 0).Distance(NewVec2(3, 4))
	if !ApproxEqual(d, 5, testEpsilon) {
		t.Errorf("expected 5, got %g", d)
	}
}

func TestVec2Lerp(t *testing.T) {
	v := NewVec2(0, 0).Lerp(NewVec2(10, 20), 0.5)
	if !v.ApproxEqual(NewVec2(5, 10), testEpsilon) {
		t.Errorf("expected (5, 10), got %v", v)
	}
}

func TestVec2Rotate(t *testing.T) {
	v := NewVec2(1, 0).Rotate(gomath.Pi / 2)
	if !v.ApproxEqual(NewVec2(0, 1), testEpsilon) {
		t.Errorf("expected (0, 1), got %v", v)
	}
}

func TestVec2Angle(t *testing.T) {
	a := NewVec2(1, 0).Angle()
	if !ApproxEqual(a, 0, testEpsilon) {
		t.Errorf("expected 0, got %g", a)
	}
	a = NewVec2(0, 1).Angle()
	if !ApproxEqual(a, gomath.Pi/2, testEpsilon) {
		t.Errorf("expected π/2, got %g", a)
	}
}

func TestVec2Reflect(t *testing.T) {
	v := NewVec2(1, -1).Reflect(NewVec2(0, 1))
	if !v.ApproxEqual(NewVec2(1, 1), testEpsilon) {
		t.Errorf("expected (1, 1), got %v", v)
	}
}

func TestVec2Perpendicular(t *testing.T) {
	p := NewVec2(1, 0).Perpendicular()
	if !p.ApproxEqual(NewVec2(0, 1), testEpsilon) {
		t.Errorf("expected (0, 1), got %v", p)
	}
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
