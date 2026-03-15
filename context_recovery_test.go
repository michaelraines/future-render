package futurerender

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewResourceTracker(t *testing.T) {
	rt := NewResourceTracker()
	require.NotNil(t, rt)
	require.Equal(t, 0, rt.ImageCount())
	require.Equal(t, 0, rt.ShaderCount())
}

func TestResourceTrackerTrackImage(t *testing.T) {
	withMockRenderer(t)
	rt := NewResourceTracker()

	img := NewImage(64, 64)
	rt.TrackImage(img, nil, true)
	require.Equal(t, 1, rt.ImageCount())

	// Tracking nil is a no-op.
	rt.TrackImage(nil, nil, false)
	require.Equal(t, 1, rt.ImageCount())
}

func TestResourceTrackerTrackImageWithPixels(t *testing.T) {
	withMockRenderer(t)
	rt := NewResourceTracker()

	pixels := make([]byte, 4*10*10)
	for i := range pixels {
		pixels[i] = 0xAB
	}

	img := &Image{width: 10, height: 10}
	rt.TrackImage(img, pixels, false)
	require.Equal(t, 1, rt.ImageCount())

	// Verify pixel data was copied (not aliased).
	pixels[0] = 0xFF
	rec := rt.images[img]
	require.Equal(t, byte(0xAB), rec.pixels[0])
}

func TestResourceTrackerUntrackImage(t *testing.T) {
	withMockRenderer(t)
	rt := NewResourceTracker()

	img := NewImage(32, 32)
	rt.TrackImage(img, nil, false)
	require.Equal(t, 1, rt.ImageCount())

	rt.UntrackImage(img)
	require.Equal(t, 0, rt.ImageCount())

	// Untracking nil is a no-op.
	rt.UntrackImage(nil)
	require.Equal(t, 0, rt.ImageCount())

	// Untracking an already-untracked image is safe.
	rt.UntrackImage(img)
	require.Equal(t, 0, rt.ImageCount())
}

func TestResourceTrackerTrackShader(t *testing.T) {
	rt := NewResourceTracker()

	s := &Shader{id: 1}
	rt.TrackShader(s, "vert", "frag", nil)
	require.Equal(t, 1, rt.ShaderCount())

	// Nil shader is a no-op.
	rt.TrackShader(nil, "v", "f", nil)
	require.Equal(t, 1, rt.ShaderCount())
}

func TestResourceTrackerUntrackShader(t *testing.T) {
	rt := NewResourceTracker()

	s := &Shader{id: 1}
	rt.TrackShader(s, "vert", "frag", nil)
	require.Equal(t, 1, rt.ShaderCount())

	rt.UntrackShader(s)
	require.Equal(t, 0, rt.ShaderCount())

	// Nil and already-untracked are safe.
	rt.UntrackShader(nil)
	rt.UntrackShader(s)
}

func TestResourceTrackerRecoverImages(t *testing.T) {
	dev, _ := withMockRenderer(t)
	rt := NewResourceTracker()

	img1 := NewImage(64, 64)
	img2 := NewImage(32, 32)
	rt.TrackImage(img1, nil, true)
	rt.TrackImage(img2, nil, false)

	initialTexCount := len(dev.textures)

	// Simulate context loss by nil-ing out textures.
	img1.texture = nil
	img1.disposed = true
	img2.texture = nil
	img2.disposed = true

	err := rt.RecoverResources(dev)
	require.NoError(t, err)

	// Textures should be recreated.
	require.Greater(t, len(dev.textures), initialTexCount)
	require.NotNil(t, img1.texture)
	require.NotNil(t, img2.texture)
	require.False(t, img1.disposed)
	require.False(t, img2.disposed)
}

func TestResourceTrackerRecoverImageWithPixels(t *testing.T) {
	dev, _ := withMockRenderer(t)
	rt := NewResourceTracker()

	pixels := make([]byte, 4*8*8)
	for i := range pixels {
		pixels[i] = 0xCD
	}

	img := &Image{width: 8, height: 8, textureID: 99}
	rt.TrackImage(img, pixels, false)

	err := rt.RecoverResources(dev)
	require.NoError(t, err)
	require.NotNil(t, img.texture)
}

func TestResourceTrackerRecoverShaders(t *testing.T) {
	dev := withShaderRenderer(t)
	rt := NewResourceTracker()

	s := &Shader{id: 42}
	rt.TrackShader(s, "vert", "frag", nil)

	s.backend = nil
	s.pipeline = nil
	s.disposed = true

	err := rt.RecoverResources(dev)
	require.NoError(t, err)
	require.NotNil(t, s.backend)
	require.NotNil(t, s.pipeline)
	require.False(t, s.disposed)
}

func TestResourceTrackerRecoverEmpty(t *testing.T) {
	dev, _ := withMockRenderer(t)
	rt := NewResourceTracker()

	err := rt.RecoverResources(dev)
	require.NoError(t, err)
}

func TestResourceTrackerRecoverNilDevice(t *testing.T) {
	rt := NewResourceTracker()

	err := rt.RecoverResources(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "device is nil")
}

func TestResourceTrackerDisposedNotRecovered(t *testing.T) {
	dev, _ := withMockRenderer(t)
	rt := NewResourceTracker()

	img := NewImage(16, 16)
	rt.TrackImage(img, nil, false)
	require.Equal(t, 1, rt.ImageCount())

	// Untrack the image (simulating Dispose calling UntrackImage).
	rt.UntrackImage(img)
	require.Equal(t, 0, rt.ImageCount())

	initialTexCount := len(dev.textures)

	err := rt.RecoverResources(dev)
	require.NoError(t, err)

	// No new textures should be created.
	require.Equal(t, initialTexCount, len(dev.textures))
}

func TestResourceTrackerRecoverRenderTarget(t *testing.T) {
	dev, _ := withMockRenderer(t)
	rt := NewResourceTracker()

	img := NewImage(128, 128)
	rt.TrackImage(img, nil, true) // with render target

	initialRTCount := len(dev.renderTargets)

	// Simulate context loss.
	img.texture = nil
	img.renderTarget = nil
	img.disposed = true

	err := rt.RecoverResources(dev)
	require.NoError(t, err)

	require.NotNil(t, img.texture)
	require.NotNil(t, img.renderTarget)
	require.Greater(t, len(dev.renderTargets), initialRTCount)
}
