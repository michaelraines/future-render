package math

import (
	"fmt"
	gomath "math"
)

// Rect represents an axis-aligned rectangle defined by minimum and maximum points.
type Rect struct {
	Min, Max Vec2
}

// NewRect creates a Rect from position and size.
func NewRect(x, y, w, h float64) Rect {
	return Rect{
		Min: Vec2{X: x, Y: y},
		Max: Vec2{X: x + w, Y: y + h},
	}
}

// RectFromMinMax creates a Rect from lo and hi corners.
func RectFromMinMax(lo, hi Vec2) Rect {
	return Rect{Min: lo, Max: hi}
}

// RectFromCenter creates a Rect centered at center with the given size.
func RectFromCenter(center Vec2, w, h float64) Rect {
	half := Vec2{X: w / 2, Y: h / 2}
	return Rect{Min: center.Sub(half), Max: center.Add(half)}
}

// Width returns the width of the rectangle.
func (r Rect) Width() float64 { return r.Max.X - r.Min.X }

// Height returns the height of the rectangle.
func (r Rect) Height() float64 { return r.Max.Y - r.Min.Y }

// Size returns the size of the rectangle as a Vec2.
func (r Rect) Size() Vec2 { return r.Max.Sub(r.Min) }

// Center returns the center of the rectangle.
func (r Rect) Center() Vec2 {
	return Vec2{
		X: (r.Min.X + r.Max.X) / 2,
		Y: (r.Min.Y + r.Max.Y) / 2,
	}
}

// Contains returns whether the point is inside the rectangle.
func (r Rect) Contains(p Vec2) bool {
	return p.X >= r.Min.X && p.X <= r.Max.X && p.Y >= r.Min.Y && p.Y <= r.Max.Y
}

// Overlaps returns whether r and other overlap.
func (r Rect) Overlaps(other Rect) bool {
	return r.Min.X < other.Max.X && r.Max.X > other.Min.X &&
		r.Min.Y < other.Max.Y && r.Max.Y > other.Min.Y
}

// Intersection returns the intersection of r and other, or a zero Rect if
// they don't overlap.
func (r Rect) Intersection(other Rect) Rect {
	result := Rect{
		Min: r.Min.Max(other.Min),
		Max: r.Max.Min(other.Max),
	}
	if result.Min.X >= result.Max.X || result.Min.Y >= result.Max.Y {
		return Rect{}
	}
	return result
}

// Union returns the smallest rectangle containing both r and other.
func (r Rect) Union(other Rect) Rect {
	return Rect{
		Min: r.Min.Min(other.Min),
		Max: r.Max.Max(other.Max),
	}
}

// Expand returns the rectangle expanded by delta in all directions.
func (r Rect) Expand(delta float64) Rect {
	d := Vec2{X: delta, Y: delta}
	return Rect{Min: r.Min.Sub(d), Max: r.Max.Add(d)}
}

// Translate returns the rectangle moved by offset.
func (r Rect) Translate(offset Vec2) Rect {
	return Rect{Min: r.Min.Add(offset), Max: r.Max.Add(offset)}
}

// IsEmpty returns whether the rectangle has zero or negative area.
func (r Rect) IsEmpty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// String returns a human-readable representation.
func (r Rect) String() string {
	return fmt.Sprintf("Rect(%g, %g, %g, %g)", r.Min.X, r.Min.Y, r.Width(), r.Height())
}

// AABB represents a 3D axis-aligned bounding box.
type AABB struct {
	Min, Max Vec3
}

// NewAABB creates an AABB from lo and hi corners.
func NewAABB(lo, hi Vec3) AABB {
	return AABB{Min: lo, Max: hi}
}

// AABBFromCenterExtents creates an AABB from a center point and half-extents.
func AABBFromCenterExtents(center, extents Vec3) AABB {
	return AABB{Min: center.Sub(extents), Max: center.Add(extents)}
}

// Center returns the center of the AABB.
func (a AABB) Center() Vec3 {
	return a.Min.Add(a.Max).Mul(0.5)
}

// Extents returns the half-extents of the AABB.
func (a AABB) Extents() Vec3 {
	return a.Max.Sub(a.Min).Mul(0.5)
}

// Size returns the full size of the AABB.
func (a AABB) Size() Vec3 {
	return a.Max.Sub(a.Min)
}

// Contains returns whether the point is inside the AABB.
func (a AABB) Contains(p Vec3) bool {
	return p.X >= a.Min.X && p.X <= a.Max.X &&
		p.Y >= a.Min.Y && p.Y <= a.Max.Y &&
		p.Z >= a.Min.Z && p.Z <= a.Max.Z
}

// Overlaps returns whether a and other overlap.
func (a AABB) Overlaps(other AABB) bool {
	return a.Min.X < other.Max.X && a.Max.X > other.Min.X &&
		a.Min.Y < other.Max.Y && a.Max.Y > other.Min.Y &&
		a.Min.Z < other.Max.Z && a.Max.Z > other.Min.Z
}

// Union returns the smallest AABB containing both a and other.
func (a AABB) Union(other AABB) AABB {
	return AABB{
		Min: a.Min.Min(other.Min),
		Max: a.Max.Max(other.Max),
	}
}

// ExpandPoint returns the AABB expanded to contain the point.
func (a AABB) ExpandPoint(p Vec3) AABB {
	return AABB{
		Min: a.Min.Min(p),
		Max: a.Max.Max(p),
	}
}

// String returns a human-readable representation.
func (a AABB) String() string {
	return fmt.Sprintf("AABB(%v, %v)", a.Min, a.Max)
}

// Frustum represents a view frustum defined by 6 planes.
// Used for frustum culling in 3D rendering.
type Frustum struct {
	Planes [6]Plane
}

// FrustumPlane indices.
const (
	FrustumLeft   = 0
	FrustumRight  = 1
	FrustumBottom = 2
	FrustumTop    = 3
	FrustumNear   = 4
	FrustumFar    = 5
)

// FrustumFromMat4 extracts a frustum from a view-projection matrix.
func FrustumFromMat4(vp Mat4) Frustum {
	var f Frustum

	// Left plane
	f.Planes[FrustumLeft] = NewPlane(
		vp[3]+vp[0], vp[7]+vp[4], vp[11]+vp[8], vp[15]+vp[12],
	).Normalize()

	// Right plane
	f.Planes[FrustumRight] = NewPlane(
		vp[3]-vp[0], vp[7]-vp[4], vp[11]-vp[8], vp[15]-vp[12],
	).Normalize()

	// Bottom plane
	f.Planes[FrustumBottom] = NewPlane(
		vp[3]+vp[1], vp[7]+vp[5], vp[11]+vp[9], vp[15]+vp[13],
	).Normalize()

	// Top plane
	f.Planes[FrustumTop] = NewPlane(
		vp[3]-vp[1], vp[7]-vp[5], vp[11]-vp[9], vp[15]-vp[13],
	).Normalize()

	// Near plane
	f.Planes[FrustumNear] = NewPlane(
		vp[3]+vp[2], vp[7]+vp[6], vp[11]+vp[10], vp[15]+vp[14],
	).Normalize()

	// Far plane
	f.Planes[FrustumFar] = NewPlane(
		vp[3]-vp[2], vp[7]-vp[6], vp[11]-vp[10], vp[15]-vp[14],
	).Normalize()

	return f
}

// ContainsPoint returns whether the frustum contains the point.
func (f Frustum) ContainsPoint(p Vec3) bool {
	for i := range f.Planes {
		if f.Planes[i].DistanceToPoint(p) < 0 {
			return false
		}
	}
	return true
}

// ContainsAABB returns whether the frustum fully or partially contains the AABB.
// Returns false only if the AABB is entirely outside any plane.
func (f Frustum) ContainsAABB(aabb AABB) bool {
	for i := range f.Planes {
		p := f.Planes[i]
		// Find the corner of the AABB most in the direction of the plane normal
		pVertex := Vec3{
			X: pickComponent(p.Normal.X >= 0, aabb.Max.X, aabb.Min.X),
			Y: pickComponent(p.Normal.Y >= 0, aabb.Max.Y, aabb.Min.Y),
			Z: pickComponent(p.Normal.Z >= 0, aabb.Max.Z, aabb.Min.Z),
		}
		if p.DistanceToPoint(pVertex) < 0 {
			return false
		}
	}
	return true
}

func pickComponent(cond bool, a, b float64) float64 {
	if cond {
		return a
	}
	return b
}

// Plane represents a plane in 3D space using the equation ax + by + cz + d = 0.
type Plane struct {
	Normal Vec3
	D      float64
}

// NewPlane creates a plane from coefficients a, b, c, d.
func NewPlane(a, b, c, d float64) Plane {
	return Plane{Normal: Vec3{X: a, Y: b, Z: c}, D: d}
}

// PlaneFromNormalPoint creates a plane from a normal and a point on the plane.
func PlaneFromNormalPoint(normal, point Vec3) Plane {
	n := normal.Normalize()
	return Plane{Normal: n, D: -n.Dot(point)}
}

// Normalize normalizes the plane equation.
func (p Plane) Normalize() Plane {
	l := p.Normal.Len()
	if l == 0 {
		return p
	}
	return Plane{Normal: p.Normal.Div(l), D: p.D / l}
}

// DistanceToPoint returns the signed distance from the plane to the point.
// Positive means the point is on the side of the normal.
func (p Plane) DistanceToPoint(point Vec3) float64 {
	return p.Normal.Dot(point) + p.D
}

// String returns a human-readable representation.
func (p Plane) String() string {
	return fmt.Sprintf("Plane(%v, %g)", p.Normal, p.D)
}

// Ray represents a ray with an origin and direction.
type Ray struct {
	Origin    Vec3
	Direction Vec3
}

// NewRay creates a ray from an origin and direction. The direction is normalized.
func NewRay(origin, direction Vec3) Ray {
	return Ray{Origin: origin, Direction: direction.Normalize()}
}

// PointAt returns the point along the ray at distance t.
func (r Ray) PointAt(t float64) Vec3 {
	return r.Origin.Add(r.Direction.Mul(t))
}

// IntersectAABB returns whether the ray intersects the AABB, and the
// distances tMin and tMax along the ray. Returns false if no intersection.
func (r Ray) IntersectAABB(aabb AABB) (tMin, tMax float64, hit bool) {
	tMin = gomath.Inf(-1)
	tMax = gomath.Inf(1)

	axes := [3][2]float64{
		{r.Origin.X, r.Direction.X},
		{r.Origin.Y, r.Direction.Y},
		{r.Origin.Z, r.Direction.Z},
	}
	mins := [3]float64{aabb.Min.X, aabb.Min.Y, aabb.Min.Z}
	maxs := [3]float64{aabb.Max.X, aabb.Max.Y, aabb.Max.Z}

	for i := 0; i < 3; i++ {
		origin, dir := axes[i][0], axes[i][1]
		if gomath.Abs(dir) < 1e-14 {
			if origin < mins[i] || origin > maxs[i] {
				return 0, 0, false
			}
			continue
		}
		invD := 1.0 / dir
		t0 := (mins[i] - origin) * invD
		t1 := (maxs[i] - origin) * invD
		if invD < 0 {
			t0, t1 = t1, t0
		}
		tMin = gomath.Max(tMin, t0)
		tMax = gomath.Min(tMax, t1)
		if tMax < tMin {
			return 0, 0, false
		}
	}
	return tMin, tMax, true
}

// IntersectPlane returns the distance along the ray to the plane intersection.
// Returns the distance and true, or 0 and false if the ray is parallel to the plane.
func (r Ray) IntersectPlane(p Plane) (float64, bool) {
	denom := p.Normal.Dot(r.Direction)
	if gomath.Abs(denom) < 1e-14 {
		return 0, false
	}
	t := -(p.Normal.Dot(r.Origin) + p.D) / denom
	return t, true
}
