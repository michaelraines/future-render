package math

import (
	gomath "math"
	"testing"
)

func TestQuatIdentity(t *testing.T) {
	q := QuatIdentity()
	v := q.RotateVec3(NewVec3(1, 2, 3))
	if !v.ApproxEqual(NewVec3(1, 2, 3), testEpsilon) {
		t.Errorf("identity should not change vector, got %v", v)
	}
}

func TestQuatFromAxisAngle(t *testing.T) {
	q := QuatFromAxisAngle(Vec3UnitZ(), gomath.Pi/2)
	v := q.RotateVec3(NewVec3(1, 0, 0))
	if !v.ApproxEqual(NewVec3(0, 1, 0), testEpsilon) {
		t.Errorf("expected (0, 1, 0), got %v", v)
	}
}

func TestQuatMul(t *testing.T) {
	q1 := QuatFromAxisAngle(Vec3UnitZ(), gomath.Pi/2)
	q2 := QuatFromAxisAngle(Vec3UnitZ(), gomath.Pi/2)
	q := q1.Mul(q2) // 180 degrees around Z
	v := q.RotateVec3(NewVec3(1, 0, 0))
	if !v.ApproxEqual(NewVec3(-1, 0, 0), testEpsilon) {
		t.Errorf("expected (-1, 0, 0), got %v", v)
	}
}

func TestQuatInverse(t *testing.T) {
	q := QuatFromAxisAngle(Vec3UnitY(), gomath.Pi/3)
	inv := q.Inverse()
	product := q.Mul(inv)
	if !product.ApproxEqual(QuatIdentity(), testEpsilon) {
		t.Errorf("q * q^-1 should be identity, got %v", product)
	}
}

func TestQuatSlerp(t *testing.T) {
	q1 := QuatIdentity()
	q2 := QuatFromAxisAngle(Vec3UnitZ(), gomath.Pi/2)
	mid := q1.Slerp(q2, 0.5)
	v := mid.RotateVec3(NewVec3(1, 0, 0))
	expected := NewVec2(1, 0).Rotate(gomath.Pi / 4)
	if !v.Vec2().ApproxEqual(expected, testEpsilon) {
		t.Errorf("expected ~(%g, %g), got (%g, %g)", expected.X, expected.Y, v.X, v.Y)
	}
}

func TestQuatToMat4(t *testing.T) {
	q := QuatFromAxisAngle(Vec3UnitZ(), gomath.Pi/2)
	m := q.ToMat4()
	v := m.MulVec3Point(NewVec3(1, 0, 0))
	if !v.ApproxEqual(NewVec3(0, 1, 0), testEpsilon) {
		t.Errorf("expected (0, 1, 0), got %v", v)
	}
}

func TestQuatToAxisAngle(t *testing.T) {
	axis := Vec3UnitY()
	angle := gomath.Pi / 3
	q := QuatFromAxisAngle(axis, angle)
	gotAxis, gotAngle := q.ToAxisAngle()
	if !gotAxis.ApproxEqual(axis, testEpsilon) {
		t.Errorf("expected axis %v, got %v", axis, gotAxis)
	}
	if !ApproxEqual(gotAngle, angle, testEpsilon) {
		t.Errorf("expected angle %g, got %g", angle, gotAngle)
	}
}

func BenchmarkQuatRotateVec3(b *testing.B) {
	q := QuatFromAxisAngle(Vec3UnitY(), 0.5)
	v := NewVec3(1, 2, 3)
	for b.Loop() {
		_ = q.RotateVec3(v)
	}
}

func BenchmarkQuatSlerp(b *testing.B) {
	q1 := QuatIdentity()
	q2 := QuatFromAxisAngle(Vec3UnitZ(), gomath.Pi/2)
	for b.Loop() {
		_ = q1.Slerp(q2, 0.5)
	}
}
