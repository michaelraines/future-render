package futurerender

import (
	"testing"

	"github.com/michaelraines/future-render/internal/backend"
	"github.com/michaelraines/future-render/internal/batch"
	fmath "github.com/michaelraines/future-render/math"
)

func TestNewImageNoRenderer(t *testing.T) {
	// Without a renderer, NewImage should still return a valid Image.
	old := globalRenderer
	globalRenderer = nil
	defer func() { globalRenderer = old }()

	img := NewImage(100, 200)
	if img == nil {
		t.Fatal("NewImage returned nil")
	}
	w, h := img.Size()
	if w != 100 || h != 200 {
		t.Errorf("expected 100x200, got %dx%d", w, h)
	}
	if img.texture != nil {
		t.Error("texture should be nil without a renderer")
	}
}

func TestSubImageUVMapping(t *testing.T) {
	img := &Image{
		width: 256, height: 256,
		textureID: 42,
		u0:        0, v0: 0, u1: 1, v1: 1,
	}

	// Sub-image: top-left quarter (0,0)-(128,128).
	sub := img.SubImage(fmath.NewRect(0, 0, 128, 128))
	if sub.width != 128 || sub.height != 128 {
		t.Errorf("expected 128x128, got %dx%d", sub.width, sub.height)
	}
	if sub.textureID != 42 {
		t.Errorf("expected textureID 42, got %d", sub.textureID)
	}
	if sub.parent != img {
		t.Error("sub-image should reference parent")
	}
	assertFloat32(t, "u0", sub.u0, 0)
	assertFloat32(t, "v0", sub.v0, 0)
	assertFloat32(t, "u1", sub.u1, 0.5)
	assertFloat32(t, "v1", sub.v1, 0.5)

	// Sub-image: bottom-right quarter (128,128)-(256,256).
	sub2 := img.SubImage(fmath.NewRect(128, 128, 128, 128))
	assertFloat32(t, "u0", sub2.u0, 0.5)
	assertFloat32(t, "v0", sub2.v0, 0.5)
	assertFloat32(t, "u1", sub2.u1, 1.0)
	assertFloat32(t, "v1", sub2.v1, 1.0)
}

func TestSubImageOfSubImage(t *testing.T) {
	root := &Image{
		width: 256, height: 256,
		textureID: 1,
		u0:        0, v0: 0, u1: 1, v1: 1,
	}

	// First sub: top-left quarter.
	sub := root.SubImage(fmath.NewRect(0, 0, 128, 128))

	// Sub of sub: top-left quarter of the sub-image.
	subsub := sub.SubImage(fmath.NewRect(0, 0, 64, 64))
	if subsub.parent != root {
		t.Error("nested sub-image should reference root parent")
	}
	assertFloat32(t, "u0", subsub.u0, 0)
	assertFloat32(t, "v0", subsub.v0, 0)
	assertFloat32(t, "u1", subsub.u1, 0.25)
	assertFloat32(t, "v1", subsub.v1, 0.25)
}

func TestDispose(t *testing.T) {
	img := NewImage(10, 10)
	img.Dispose()
	if !img.disposed {
		t.Error("image should be disposed")
	}

	// DrawImage on disposed image should be a no-op.
	img.DrawImage(NewImage(5, 5), nil) // should not panic
}

func TestSubImageDisposeDoesNotReleaseParent(t *testing.T) {
	root := &Image{
		width: 64, height: 64,
		textureID: 1,
		u0:        0, v0: 0, u1: 1, v1: 1,
	}
	sub := root.SubImage(fmath.NewRect(0, 0, 32, 32))
	sub.Dispose()
	if root.disposed {
		t.Error("disposing sub-image should not dispose root")
	}
}

func TestDrawImageSubmitsToBatcher(t *testing.T) {
	b := batch.NewBatcher(1024, 1024)
	rend := &renderer{
		batcher:        b,
		whiteTextureID: 1,
	}
	old := globalRenderer
	globalRenderer = rend
	defer func() { globalRenderer = old }()

	dst := &Image{
		width: 320, height: 240,
		u0: 0, v0: 0, u1: 1, v1: 1,
	}
	src := &Image{
		width: 64, height: 64,
		textureID: 5,
		u0:        0, v0: 0, u1: 1, v1: 1,
	}

	opts := &DrawImageOptions{}
	opts.GeoM.Translate(100, 50)
	dst.DrawImage(src, opts)

	if b.CommandCount() != 1 {
		t.Fatalf("expected 1 command, got %d", b.CommandCount())
	}

	// Flush and verify batch contents.
	batches := b.Flush()
	if len(batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}
	got := batches[0]
	if len(got.Vertices) != 4 {
		t.Errorf("expected 4 vertices, got %d", len(got.Vertices))
	}
	if len(got.Indices) != 6 {
		t.Errorf("expected 6 indices, got %d", len(got.Indices))
	}

	// Verify the translated position of first vertex (top-left).
	v0 := got.Vertices[0]
	assertFloat32(t, "v0.PosX", v0.PosX, 100)
	assertFloat32(t, "v0.PosY", v0.PosY, 50)

	// Verify the bottom-right vertex.
	v2 := got.Vertices[2]
	assertFloat32(t, "v2.PosX", v2.PosX, 164)
	assertFloat32(t, "v2.PosY", v2.PosY, 114)
}

func TestDrawImageColorScale(t *testing.T) {
	b := batch.NewBatcher(1024, 1024)
	rend := &renderer{
		batcher:        b,
		whiteTextureID: 1,
	}
	old := globalRenderer
	globalRenderer = rend
	defer func() { globalRenderer = old }()

	dst := &Image{width: 100, height: 100, u0: 0, v0: 0, u1: 1, v1: 1}
	src := &Image{width: 32, height: 32, textureID: 2, u0: 0, v0: 0, u1: 1, v1: 1}

	opts := &DrawImageOptions{
		ColorScale: fmath.Color{R: 0.5, G: 0.5, B: 0.5, A: 0.5},
	}
	dst.DrawImage(src, opts)

	batches := b.Flush()
	v := batches[0].Vertices[0]
	assertFloat32(t, "R", v.R, 0.5)
	assertFloat32(t, "G", v.G, 0.5)
	assertFloat32(t, "B", v.B, 0.5)
	assertFloat32(t, "A", v.A, 0.5)
}

func TestDrawImageDefaultColorIsWhite(t *testing.T) {
	b := batch.NewBatcher(1024, 1024)
	rend := &renderer{
		batcher:        b,
		whiteTextureID: 1,
	}
	old := globalRenderer
	globalRenderer = rend
	defer func() { globalRenderer = old }()

	dst := &Image{width: 100, height: 100, u0: 0, v0: 0, u1: 1, v1: 1}
	src := &Image{width: 32, height: 32, textureID: 2, u0: 0, v0: 0, u1: 1, v1: 1}

	dst.DrawImage(src, nil) // nil opts → default color

	batches := b.Flush()
	v := batches[0].Vertices[0]
	assertFloat32(t, "R", v.R, 1)
	assertFloat32(t, "G", v.G, 1)
	assertFloat32(t, "B", v.B, 1)
	assertFloat32(t, "A", v.A, 1)
}

func TestFillSubmitsToBatcher(t *testing.T) {
	b := batch.NewBatcher(1024, 1024)
	rend := &renderer{
		batcher:        b,
		whiteTextureID: 99,
	}
	old := globalRenderer
	globalRenderer = rend
	defer func() { globalRenderer = old }()

	img := &Image{width: 320, height: 240, u0: 0, v0: 0, u1: 1, v1: 1}
	img.Fill(fmath.Color{R: 1, G: 0, B: 0, A: 1})

	batches := b.Flush()
	if len(batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}
	if batches[0].TextureID != 99 {
		t.Errorf("expected white texture ID 99, got %d", batches[0].TextureID)
	}
	v := batches[0].Vertices[0]
	assertFloat32(t, "R", v.R, 1)
	assertFloat32(t, "G", v.G, 0)
}

func TestBlendToBackend(t *testing.T) {
	tests := []struct {
		pub  BlendMode
		want backend.BlendMode
	}{
		{BlendSourceOver, backend.BlendSourceOver},
		{BlendAdditive, backend.BlendAdditive},
		{BlendMultiplicative, backend.BlendMultiplicative},
	}
	for _, tt := range tests {
		got := blendToBackend(tt.pub)
		if got != tt.want {
			t.Errorf("blendToBackend(%d) = %d, want %d", tt.pub, got, tt.want)
		}
	}
}

func assertFloat32(t *testing.T, name string, got, want float32) {
	t.Helper()
	const eps = 1e-5
	if diff := got - want; diff > eps || diff < -eps {
		t.Errorf("%s: got %g, want %g", name, got, want)
	}
}
