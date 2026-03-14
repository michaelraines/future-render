package futurerender

import (
	fmath "github.com/michaelraines/future-render/math"
)

// Image represents a renderable image (texture). It can be used as a
// render target or as a source for drawing operations.
//
// This type is the equivalent of ebiten.Image.
type Image struct {
	width, height int
	disposed      bool
}

// NewImage creates a new blank image with the given dimensions.
func NewImage(width, height int) *Image {
	return &Image{
		width:  width,
		height: height,
	}
}

// Size returns the image dimensions.
func (img *Image) Size() (width, height int) {
	return img.width, img.height
}

// Bounds returns the image bounds as a Rect.
func (img *Image) Bounds() fmath.Rect {
	return fmath.NewRect(0, 0, float64(img.width), float64(img.height))
}

// DrawImage draws src onto img with the given options.
func (img *Image) DrawImage(src *Image, opts *DrawImageOptions) {
	if img.disposed || src == nil || src.disposed {
		return
	}
	// Implementation will submit a draw command to the batcher via
	// the current frame's command list. Placeholder for API shape.
}

// DrawTriangles draws triangles with the given vertices, indices, and options.
// This is the low-level drawing primitive equivalent to ebiten.DrawTriangles.
func (img *Image) DrawTriangles(vertices []Vertex, indices []uint16, src *Image, opts *DrawTrianglesOptions) {
	if img.disposed {
		return
	}
	// Implementation will submit vertex/index data to the batcher.
}

// Fill fills the entire image with the given color.
func (img *Image) Fill(c fmath.Color) {
	if img.disposed {
		return
	}
	// Implementation will issue a clear command.
}

// SubImage returns a sub-region of the image for sprite sheet support.
func (img *Image) SubImage(r fmath.Rect) *Image {
	// In the full implementation, this returns an Image that references
	// the same underlying texture but with UV coordinates mapped to the
	// sub-region.
	return &Image{
		width:  int(r.Width()),
		height: int(r.Height()),
	}
}

// Dispose releases the image's GPU resources.
func (img *Image) Dispose() {
	img.disposed = true
}

// DrawImageOptions holds options for DrawImage.
type DrawImageOptions struct {
	// GeoM is the geometry transformation matrix (2D affine transform).
	GeoM GeoM

	// ColorScale scales the RGBA color of each pixel.
	ColorScale fmath.Color

	// ColorM is the color matrix transformation.
	ColorM fmath.ColorMatrix

	// Blend specifies the blend mode.
	Blend BlendMode

	// Filter specifies the texture filter.
	Filter Filter
}

// DrawTrianglesOptions holds options for DrawTriangles.
type DrawTrianglesOptions struct {
	// Blend specifies the blend mode.
	Blend BlendMode

	// Filter specifies the texture filter.
	Filter Filter

	// FillRule specifies the fill rule for overlapping triangles.
	FillRule FillRule
}

// Vertex represents a vertex for DrawTriangles.
type Vertex struct {
	DstX, DstY                     float32
	SrcX, SrcY                     float32
	ColorR, ColorG, ColorB, ColorA float32
}

// GeoM represents a 2D affine transformation matrix.
// This provides an API compatible with ebiten.GeoM.
type GeoM struct {
	m fmath.Mat3
}

// NewGeoM creates an identity GeoM.
func NewGeoM() GeoM {
	return GeoM{m: fmath.Mat3Identity()}
}

// Translate adds a translation to the transformation.
func (g *GeoM) Translate(tx, ty float64) {
	g.m = fmath.Mat3Translate(tx, ty).Mul(g.m)
}

// Scale adds a scaling to the transformation.
func (g *GeoM) Scale(sx, sy float64) {
	g.m = fmath.Mat3Scale(sx, sy).Mul(g.m)
}

// Rotate adds a rotation (radians) to the transformation.
func (g *GeoM) Rotate(angle float64) {
	g.m = fmath.Mat3Rotate(angle).Mul(g.m)
}

// Skew adds a shear/skew to the transformation.
func (g *GeoM) Skew(sx, sy float64) {
	g.m = fmath.Mat3Shear(sx, sy).Mul(g.m)
}

// Concat concatenates another GeoM onto this one.
func (g *GeoM) Concat(other GeoM) {
	g.m = other.m.Mul(g.m)
}

// Reset resets the GeoM to identity.
func (g *GeoM) Reset() {
	g.m = fmath.Mat3Identity()
}

// Apply transforms a point by this GeoM.
func (g *GeoM) Apply(x, y float64) (rx, ry float64) {
	v := g.m.MulVec2(fmath.NewVec2(x, y))
	return v.X, v.Y
}

// Mat3 returns the underlying 3x3 matrix.
func (g *GeoM) Mat3() fmath.Mat3 {
	return g.m
}

// BlendMode specifies how colors are blended.
type BlendMode int

// BlendMode constants.
const (
	BlendSourceOver     BlendMode = iota // Standard alpha blending
	BlendAdditive                        // Additive blending
	BlendMultiplicative                  // Multiplicative blending
)

// Filter specifies texture filtering.
type Filter int

// Filter constants.
const (
	FilterNearest Filter = iota
	FilterLinear
)

// FillRule specifies the fill rule for overlapping triangles.
type FillRule int

// FillRule constants.
const (
	FillRuleNonZero FillRule = iota
	FillRuleEvenOdd
)
