package futurerender

import (
	"sync/atomic"

	"github.com/michaelraines/future-render/internal/backend"
	"github.com/michaelraines/future-render/internal/batch"
)

// renderer holds internal rendering state shared between the engine loop
// and public API types like Image. It is initialized by the platform engine.
type renderer struct {
	device  backend.Device
	batcher *batch.Batcher

	// Monotonic texture ID counter for batcher sort keys.
	nextTextureID atomic.Uint32

	// whiteTextureID is the texture ID for a 1x1 white texture,
	// used for untextured draws (Fill, etc.).
	whiteTextureID uint32

	// registerTexture is called when a new texture is created (e.g. by
	// NewImage) so the engine can track it for lookup during rendering.
	registerTexture func(id uint32, tex backend.Texture)
}

// globalRenderer is the active renderer, set by the engine during init.
var globalRenderer *renderer

func (r *renderer) allocTextureID() uint32 {
	return r.nextTextureID.Add(1)
}
