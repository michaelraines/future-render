package math

import (
	gomath "math"
	"testing"
)

func TestMat4Identity(t *testing.T) {
	m := Mat4Identity()
	v := m.MulVec4(NewVec4(1, 2, 3, 1))
	if !v.ApproxEqual(NewVec4(1, 2, 3, 1), testEpsilon) {
		t.Errorf("identity should not change vector, got %v", v)
	}
}

func TestMat4Translate(t *testing.T) {
	m := Mat4Translate(10, 20, 30)
	v := m.MulVec3Point(NewVec3(1, 2, 3))
	if !v.ApproxEqual(NewVec3(11, 22, 33), testEpsilon) {
		t.Errorf("expected (11, 22, 33), got %v", v)
	}
}

func TestMat4Scale(t *testing.T) {
	m := Mat4Scale(2, 3, 4)
	v := m.MulVec3Point(NewVec3(1, 1, 1))
	if !v.ApproxEqual(NewVec3(2, 3, 4), testEpsilon) {
		t.Errorf("expected (2, 3, 4), got %v", v)
	}
}

func TestMat4RotateZ(t *testing.T) {
	m := Mat4RotateZ(gomath.Pi / 2)
	v := m.MulVec3Point(NewVec3(1, 0, 0))
	if !v.ApproxEqual(NewVec3(0, 1, 0), testEpsilon) {
		t.Errorf("expected (0, 1, 0), got %v", v)
	}
}

func TestMat4MulIdentity(t *testing.T) {
	a := Mat4Translate(1, 2, 3)
	b := Mat4Identity()
	c := a.Mul(b)
	if !c.ApproxEqual(a, testEpsilon) {
		t.Errorf("M * I should equal M")
	}
}

func TestMat4Inverse(t *testing.T) {
	m := Mat4Translate(5, 10, 15).Mul(Mat4Scale(2, 3, 4))
	inv, ok := m.Inverse()
	if !ok {
		t.Fatal("expected invertible matrix")
	}
	product := m.Mul(inv)
	if !product.ApproxEqual(Mat4Identity(), 1e-10) {
		t.Errorf("M * M^-1 should be identity, got %v", product)
	}
}

func TestMat4InverseSingular(t *testing.T) {
	m := Mat4{} // zero matrix
	_, ok := m.Inverse()
	if ok {
		t.Error("zero matrix should not be invertible")
	}
}

func TestMat4Determinant(t *testing.T) {
	d := Mat4Identity().Determinant()
	if !ApproxEqual(d, 1, testEpsilon) {
		t.Errorf("identity determinant should be 1, got %g", d)
	}
	d = Mat4Scale(2, 3, 4).Determinant()
	if !ApproxEqual(d, 24, testEpsilon) {
		t.Errorf("expected 24, got %g", d)
	}
}

func TestMat4Ortho(t *testing.T) {
	m := Mat4Ortho(0, 800, 600, 0, -1, 1)
	// Center of the screen should map to (0, 0, 0)
	center := m.MulVec3Point(NewVec3(400, 300, 0))
	if !center.ApproxEqual(Vec3Zero(), testEpsilon) {
		t.Errorf("screen center should map to origin, got %v", center)
	}
}

func TestMat4Perspective(t *testing.T) {
	m := Mat4Perspective(gomath.Pi/4, 16.0/9.0, 0.1, 100)
	// A point at the near plane center should map to (0, 0, -1) in NDC
	p := m.MulVec4(NewVec4(0, 0, -0.1, 1))
	ndc := p.PerspectiveDivide()
	if !ApproxEqual(ndc.Z, -1, 1e-6) {
		t.Errorf("near plane z should be -1 in NDC, got %g", ndc.Z)
	}
}

func TestMat4LookAt(t *testing.T) {
	m := Mat4LookAt(NewVec3(0, 0, 5), Vec3Zero(), Vec3UnitY())
	// The eye position should map to origin
	p := m.MulVec3Point(NewVec3(0, 0, 5))
	if !p.ApproxEqual(Vec3Zero(), testEpsilon) {
		t.Errorf("eye should map to origin, got %v", p)
	}
}

func TestMat4Translation(t *testing.T) {
	m := Mat4Translate(1, 2, 3)
	tr := m.Translation()
	if !tr.ApproxEqual(NewVec3(1, 2, 3), testEpsilon) {
		t.Errorf("expected (1, 2, 3), got %v", tr)
	}
}

func BenchmarkMat4Mul(b *testing.B) {
	a := Mat4Translate(1, 2, 3)
	c := Mat4RotateZ(0.5)
	for b.Loop() {
		_ = a.Mul(c)
	}
}

func BenchmarkMat4Inverse(b *testing.B) {
	m := Mat4Translate(1, 2, 3).Mul(Mat4Scale(2, 3, 4))
	for b.Loop() {
		_, _ = m.Inverse()
	}
}

func BenchmarkMat4MulVec4(b *testing.B) {
	m := Mat4Translate(1, 2, 3)
	v := NewVec4(1, 2, 3, 1)
	for b.Loop() {
		_ = m.MulVec4(v)
	}
}
