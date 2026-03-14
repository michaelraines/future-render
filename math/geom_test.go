package math

import (
	gomath "math"
	"testing"
)

func TestRectContains(t *testing.T) {
	r := NewRect(0, 0, 100, 100)
	if !r.Contains(NewVec2(50, 50)) {
		t.Error("should contain center")
	}
	if r.Contains(NewVec2(101, 50)) {
		t.Error("should not contain point outside")
	}
}

func TestRectOverlaps(t *testing.T) {
	a := NewRect(0, 0, 10, 10)
	b := NewRect(5, 5, 10, 10)
	if !a.Overlaps(b) {
		t.Error("rectangles should overlap")
	}
	c := NewRect(20, 20, 5, 5)
	if a.Overlaps(c) {
		t.Error("rectangles should not overlap")
	}
}

func TestRectIntersection(t *testing.T) {
	a := NewRect(0, 0, 10, 10)
	b := NewRect(5, 5, 10, 10)
	inter := a.Intersection(b)
	expected := NewRect(5, 5, 5, 5)
	if inter != expected {
		t.Errorf("expected %v, got %v", expected, inter)
	}
}

func TestRectUnion(t *testing.T) {
	a := NewRect(0, 0, 5, 5)
	b := NewRect(3, 3, 5, 5)
	u := a.Union(b)
	expected := NewRect(0, 0, 8, 8)
	if u != expected {
		t.Errorf("expected %v, got %v", expected, u)
	}
}

func TestAABBContains(t *testing.T) {
	aabb := NewAABB(Vec3Zero(), Vec3One())
	if !aabb.Contains(NewVec3(0.5, 0.5, 0.5)) {
		t.Error("should contain center")
	}
	if aabb.Contains(NewVec3(2, 0, 0)) {
		t.Error("should not contain outside point")
	}
}

func TestRayIntersectAABB(t *testing.T) {
	aabb := NewAABB(NewVec3(-1, -1, -1), Vec3One())
	ray := NewRay(NewVec3(0, 0, 5), NewVec3(0, 0, -1))
	tMin, tMax, hit := ray.IntersectAABB(aabb)
	if !hit {
		t.Fatal("ray should hit AABB")
	}
	if !ApproxEqual(tMin, 4, testEpsilon) || !ApproxEqual(tMax, 6, testEpsilon) {
		t.Errorf("expected tMin=4, tMax=6, got %g, %g", tMin, tMax)
	}
}

func TestRayMissAABB(t *testing.T) {
	aabb := NewAABB(NewVec3(-1, -1, -1), Vec3One())
	ray := NewRay(NewVec3(5, 5, 5), NewVec3(0, 0, -1))
	_, _, hit := ray.IntersectAABB(aabb)
	if hit {
		t.Error("ray should miss AABB")
	}
}

func TestFrustumContainsPoint(t *testing.T) {
	vp := Mat4Perspective(gomath.Pi/4, 1.0, 0.1, 100).Mul(
		Mat4LookAt(NewVec3(0, 0, 5), Vec3Zero(), Vec3UnitY()),
	)
	f := FrustumFromMat4(vp)
	if !f.ContainsPoint(Vec3Zero()) {
		t.Error("frustum should contain origin")
	}
}

func TestPlaneDistance(t *testing.T) {
	p := PlaneFromNormalPoint(Vec3UnitY(), Vec3Zero())
	d := p.DistanceToPoint(NewVec3(0, 5, 0))
	if !ApproxEqual(d, 5, testEpsilon) {
		t.Errorf("expected distance 5, got %g", d)
	}
}
