package math

import (
	gomath "math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRectContains(t *testing.T) {
	r := NewRect(0, 0, 100, 100)
	require.True(t, r.Contains(NewVec2(50, 50)))
	require.False(t, r.Contains(NewVec2(101, 50)))
}

func TestRectOverlaps(t *testing.T) {
	a := NewRect(0, 0, 10, 10)
	b := NewRect(5, 5, 10, 10)
	require.True(t, a.Overlaps(b))

	c := NewRect(20, 20, 5, 5)
	require.False(t, a.Overlaps(c))
}

func TestRectIntersection(t *testing.T) {
	a := NewRect(0, 0, 10, 10)
	b := NewRect(5, 5, 10, 10)
	inter := a.Intersection(b)
	expected := NewRect(5, 5, 5, 5)
	require.Equal(t, expected, inter)

	// Non-overlapping produces zero rect
	c := NewRect(20, 20, 5, 5)
	require.Equal(t, Rect{}, a.Intersection(c))
}

func TestRectUnion(t *testing.T) {
	a := NewRect(0, 0, 5, 5)
	b := NewRect(3, 3, 5, 5)
	u := a.Union(b)
	expected := NewRect(0, 0, 8, 8)
	require.Equal(t, expected, u)
}

func TestAABBContains(t *testing.T) {
	aabb := NewAABB(Vec3Zero(), Vec3One())
	require.True(t, aabb.Contains(NewVec3(0.5, 0.5, 0.5)))
	require.False(t, aabb.Contains(NewVec3(2, 0, 0)))
}

func TestRayIntersectAABB(t *testing.T) {
	aabb := NewAABB(NewVec3(-1, -1, -1), Vec3One())
	ray := NewRay(NewVec3(0, 0, 5), NewVec3(0, 0, -1))
	tMin, tMax, hit := ray.IntersectAABB(aabb)
	require.True(t, hit, "ray should hit AABB")
	require.InDelta(t, 4.0, tMin, 1e-9)
	require.InDelta(t, 6.0, tMax, 1e-9)
}

func TestRayMissAABB(t *testing.T) {
	aabb := NewAABB(NewVec3(-1, -1, -1), Vec3One())
	ray := NewRay(NewVec3(5, 5, 5), NewVec3(0, 0, -1))
	_, _, hit := ray.IntersectAABB(aabb)
	require.False(t, hit, "ray should miss AABB")
}

func TestFrustumContainsPoint(t *testing.T) {
	vp := Mat4Perspective(gomath.Pi/4, 1.0, 0.1, 100).Mul(
		Mat4LookAt(NewVec3(0, 0, 5), Vec3Zero(), Vec3UnitY()),
	)
	f := FrustumFromMat4(vp)
	require.True(t, f.ContainsPoint(Vec3Zero()))
}

func TestPlaneDistance(t *testing.T) {
	p := PlaneFromNormalPoint(Vec3UnitY(), Vec3Zero())
	d := p.DistanceToPoint(NewVec3(0, 5, 0))
	require.InDelta(t, 5.0, d, 1e-9)
}

// --- New Rect tests ---

func TestRectWidthHeight(t *testing.T) {
	r := NewRect(10, 20, 30, 40)
	require.InDelta(t, 30.0, r.Width(), 1e-9)
	require.InDelta(t, 40.0, r.Height(), 1e-9)
}

func TestRectCenter(t *testing.T) {
	r := NewRect(0, 0, 100, 200)
	c := r.Center()
	require.InDelta(t, 50.0, c.X, 1e-9)
	require.InDelta(t, 100.0, c.Y, 1e-9)
}

func TestRectSize(t *testing.T) {
	r := NewRect(5, 10, 20, 30)
	s := r.Size()
	require.InDelta(t, 20.0, s.X, 1e-9)
	require.InDelta(t, 30.0, s.Y, 1e-9)
}

func TestRectExpand(t *testing.T) {
	r := NewRect(10, 10, 20, 20)
	e := r.Expand(5)
	require.InDelta(t, 5.0, e.Min.X, 1e-9)
	require.InDelta(t, 5.0, e.Min.Y, 1e-9)
	require.InDelta(t, 35.0, e.Max.X, 1e-9)
	require.InDelta(t, 35.0, e.Max.Y, 1e-9)
}

func TestRectTranslate(t *testing.T) {
	r := NewRect(0, 0, 10, 10)
	moved := r.Translate(NewVec2(5, 3))
	require.InDelta(t, 5.0, moved.Min.X, 1e-9)
	require.InDelta(t, 3.0, moved.Min.Y, 1e-9)
	require.InDelta(t, 15.0, moved.Max.X, 1e-9)
	require.InDelta(t, 13.0, moved.Max.Y, 1e-9)
}

func TestRectIsEmpty(t *testing.T) {
	require.True(t, Rect{}.IsEmpty())
	require.False(t, NewRect(0, 0, 10, 10).IsEmpty())
	// Negative area
	require.True(t, RectFromMinMax(NewVec2(10, 10), NewVec2(5, 5)).IsEmpty())
}

func TestRectString(t *testing.T) {
	s := NewRect(1, 2, 3, 4).String()
	require.Equal(t, "Rect(1, 2, 3, 4)", s)
}

func TestRectFromMinMax(t *testing.T) {
	r := RectFromMinMax(NewVec2(1, 2), NewVec2(11, 12))
	require.InDelta(t, 10.0, r.Width(), 1e-9)
	require.InDelta(t, 10.0, r.Height(), 1e-9)
}

func TestRectFromCenter(t *testing.T) {
	r := RectFromCenter(NewVec2(50, 50), 20, 30)
	require.InDelta(t, 40.0, r.Min.X, 1e-9)
	require.InDelta(t, 35.0, r.Min.Y, 1e-9)
	require.InDelta(t, 60.0, r.Max.X, 1e-9)
	require.InDelta(t, 65.0, r.Max.Y, 1e-9)
}

// --- New AABB tests ---

func TestAABBCenter(t *testing.T) {
	a := NewAABB(NewVec3(0, 0, 0), NewVec3(10, 20, 30))
	c := a.Center()
	require.InDelta(t, 5.0, c.X, 1e-9)
	require.InDelta(t, 10.0, c.Y, 1e-9)
	require.InDelta(t, 15.0, c.Z, 1e-9)
}

func TestAABBExtents(t *testing.T) {
	a := NewAABB(NewVec3(0, 0, 0), NewVec3(10, 20, 30))
	e := a.Extents()
	require.InDelta(t, 5.0, e.X, 1e-9)
	require.InDelta(t, 10.0, e.Y, 1e-9)
	require.InDelta(t, 15.0, e.Z, 1e-9)
}

func TestAABBSize(t *testing.T) {
	a := NewAABB(NewVec3(1, 2, 3), NewVec3(4, 6, 9))
	s := a.Size()
	require.InDelta(t, 3.0, s.X, 1e-9)
	require.InDelta(t, 4.0, s.Y, 1e-9)
	require.InDelta(t, 6.0, s.Z, 1e-9)
}

func TestAABBOverlaps(t *testing.T) {
	a := NewAABB(Vec3Zero(), Vec3One())
	b := NewAABB(NewVec3(0.5, 0.5, 0.5), NewVec3(2, 2, 2))
	require.True(t, a.Overlaps(b))

	c := NewAABB(NewVec3(5, 5, 5), NewVec3(6, 6, 6))
	require.False(t, a.Overlaps(c))
}

func TestAABBUnion(t *testing.T) {
	a := NewAABB(Vec3Zero(), Vec3One())
	b := NewAABB(NewVec3(2, 2, 2), NewVec3(3, 3, 3))
	u := a.Union(b)
	require.InDelta(t, 0.0, u.Min.X, 1e-9)
	require.InDelta(t, 3.0, u.Max.X, 1e-9)
}

func TestAABBExpandPoint(t *testing.T) {
	a := NewAABB(Vec3Zero(), Vec3One())
	expanded := a.ExpandPoint(NewVec3(5, -1, 0.5))
	require.InDelta(t, 0.0, expanded.Min.X, 1e-9)
	require.InDelta(t, -1.0, expanded.Min.Y, 1e-9)
	require.InDelta(t, 5.0, expanded.Max.X, 1e-9)
}

func TestAABBFromCenterExtents(t *testing.T) {
	a := AABBFromCenterExtents(NewVec3(5, 5, 5), NewVec3(2, 2, 2))
	require.InDelta(t, 3.0, a.Min.X, 1e-9)
	require.InDelta(t, 7.0, a.Max.X, 1e-9)
}

func TestAABBString(t *testing.T) {
	a := NewAABB(Vec3Zero(), Vec3One())
	s := a.String()
	require.Contains(t, s, "AABB(")
}

// --- Frustum AABB test ---

func TestFrustumContainsAABB(t *testing.T) {
	vp := Mat4Perspective(gomath.Pi/4, 1.0, 0.1, 100).Mul(
		Mat4LookAt(NewVec3(0, 0, 5), Vec3Zero(), Vec3UnitY()),
	)
	f := FrustumFromMat4(vp)

	// Small box at origin should be inside
	inside := NewAABB(NewVec3(-0.1, -0.1, -0.1), NewVec3(0.1, 0.1, 0.1))
	require.True(t, f.ContainsAABB(inside))

	// Large box far away should be outside
	outside := NewAABB(NewVec3(1000, 1000, 1000), NewVec3(1001, 1001, 1001))
	require.False(t, f.ContainsAABB(outside))
}

// --- Ray-Plane intersection test ---

func TestRayIntersectPlane(t *testing.T) {
	p := PlaneFromNormalPoint(Vec3UnitY(), Vec3Zero())
	ray := NewRay(NewVec3(0, 5, 0), NewVec3(0, -1, 0))
	dist, hit := ray.IntersectPlane(p)
	require.True(t, hit)
	require.InDelta(t, 5.0, dist, 1e-9)
}

func TestRayIntersectPlaneParallel(t *testing.T) {
	p := PlaneFromNormalPoint(Vec3UnitY(), Vec3Zero())
	ray := NewRay(NewVec3(0, 5, 0), NewVec3(1, 0, 0))
	_, hit := ray.IntersectPlane(p)
	require.False(t, hit)
}

// --- Plane tests ---

func TestPlaneNormalize(t *testing.T) {
	p := NewPlane(0, 3, 4, 10)
	pn := p.Normalize()
	require.InDelta(t, 1.0, pn.Normal.Len(), 1e-9)
}

func TestPlaneNormalizeZero(t *testing.T) {
	p := NewPlane(0, 0, 0, 0)
	pn := p.Normalize()
	require.InDelta(t, 0.0, pn.Normal.Len(), 1e-9)
}

func TestPlaneString(t *testing.T) {
	p := PlaneFromNormalPoint(Vec3UnitY(), Vec3Zero())
	s := p.String()
	require.Contains(t, s, "Plane(")
}

// --- Ray PointAt test ---

func TestRayPointAt(t *testing.T) {
	ray := NewRay(NewVec3(1, 2, 3), NewVec3(1, 0, 0))
	p := ray.PointAt(5)
	require.InDelta(t, 6.0, p.X, 1e-9)
	require.InDelta(t, 2.0, p.Y, 1e-9)
	require.InDelta(t, 3.0, p.Z, 1e-9)
}

// --- Test ray with axis-aligned direction (parallel to an AABB face) ---

func TestRayIntersectAABBParallelAxis(t *testing.T) {
	aabb := NewAABB(NewVec3(-1, -1, -1), NewVec3(1, 1, 1))
	// Ray parallel to X axis, inside Y and Z range
	ray := NewRay(NewVec3(-5, 0, 0), NewVec3(1, 0, 0))
	tMin, tMax, hit := ray.IntersectAABB(aabb)
	require.True(t, hit)
	require.InDelta(t, 4.0, tMin, 1e-9)
	require.InDelta(t, 6.0, tMax, 1e-9)

	// Ray parallel to X axis, outside Y range
	ray2 := NewRay(NewVec3(-5, 5, 0), NewVec3(1, 0, 0))
	_, _, hit2 := ray2.IntersectAABB(aabb)
	require.False(t, hit2)
}

func TestFrustumContainsPointOutside(t *testing.T) {
	vp := Mat4Perspective(gomath.Pi/4, 1.0, 0.1, 100).Mul(
		Mat4LookAt(NewVec3(0, 0, 5), Vec3Zero(), Vec3UnitY()),
	)
	f := FrustumFromMat4(vp)
	require.False(t, f.ContainsPoint(NewVec3(1000, 0, 0)))
}
