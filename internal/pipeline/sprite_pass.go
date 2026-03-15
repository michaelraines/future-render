package pipeline

import (
	"unsafe"

	"github.com/michaelraines/future-render/internal/backend"
	"github.com/michaelraines/future-render/internal/batch"
)

// TextureResolver maps a batcher texture ID to a backend.Texture.
type TextureResolver func(textureID uint32) backend.Texture

// SpritePass renders 2D sprite batches. It flushes the batcher, uploads
// vertex/index data to dynamic GPU buffers, and issues indexed draw calls.
type SpritePass struct {
	batcher  *batch.Batcher
	pipeline backend.Pipeline
	shader   backend.Shader

	// Dynamic GPU buffers for per-frame vertex/index uploads.
	vertexBuf backend.Buffer
	indexBuf  backend.Buffer

	// ResolveTexture maps batch texture IDs to backend textures.
	ResolveTexture TextureResolver

	// Projection is the orthographic projection matrix, set each frame.
	Projection [16]float32
}

// SpritePassConfig holds configuration for creating a SpritePass.
type SpritePassConfig struct {
	Device   backend.Device
	Batcher  *batch.Batcher
	Pipeline backend.Pipeline
	Shader   backend.Shader

	// MaxVertices is the capacity of the dynamic vertex buffer.
	MaxVertices int

	// MaxIndices is the capacity of the dynamic index buffer.
	MaxIndices int
}

// NewSpritePass creates a new sprite pass with pre-allocated GPU buffers.
func NewSpritePass(cfg SpritePassConfig) (*SpritePass, error) {
	vbuf, err := cfg.Device.NewBuffer(backend.BufferDescriptor{
		Size:    cfg.MaxVertices * batch.Vertex2DSize,
		Usage:   backend.BufferUsageVertex,
		Dynamic: true,
	})
	if err != nil {
		return nil, err
	}

	ibuf, err := cfg.Device.NewBuffer(backend.BufferDescriptor{
		Size:    cfg.MaxIndices * 2, // uint16 indices
		Usage:   backend.BufferUsageIndex,
		Dynamic: true,
	})
	if err != nil {
		vbuf.Dispose()
		return nil, err
	}

	return &SpritePass{
		batcher:   cfg.Batcher,
		pipeline:  cfg.Pipeline,
		shader:    cfg.Shader,
		vertexBuf: vbuf,
		indexBuf:  ibuf,
	}, nil
}

// Name returns the pass name.
func (sp *SpritePass) Name() string { return "sprite" }

// Execute flushes the batcher and renders all batches.
func (sp *SpritePass) Execute(enc backend.CommandEncoder, ctx *PassContext) {
	batches := sp.batcher.Flush()
	if len(batches) == 0 {
		return
	}

	// Set pipeline and projection uniform.
	enc.SetPipeline(sp.pipeline)
	sp.shader.SetUniformMat4("uProjection", sp.Projection)
	sp.shader.SetUniformInt("uTexture", 0)

	for i := range batches {
		b := &batches[i]

		// Upload vertex data.
		vertexData := vertexSliceToBytes(b.Vertices)
		sp.vertexBuf.Upload(vertexData)
		enc.SetVertexBuffer(sp.vertexBuf, 0)

		// Upload index data.
		indexData := indexSliceToBytes(b.Indices)
		sp.indexBuf.Upload(indexData)
		enc.SetIndexBuffer(sp.indexBuf, backend.IndexUint16)

		// Bind texture.
		if sp.ResolveTexture != nil {
			tex := sp.ResolveTexture(b.TextureID)
			if tex != nil {
				enc.SetTexture(tex, 0)
			}
		}

		// Draw.
		enc.DrawIndexed(len(b.Indices), 1, 0)
	}
}

// Dispose releases the pass's GPU buffers.
func (sp *SpritePass) Dispose() {
	if sp.vertexBuf != nil {
		sp.vertexBuf.Dispose()
	}
	if sp.indexBuf != nil {
		sp.indexBuf.Dispose()
	}
}

// vertexSliceToBytes reinterprets a []Vertex2D as a []byte without copying.
func vertexSliceToBytes(verts []batch.Vertex2D) []byte {
	if len(verts) == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(&verts[0])), len(verts)*batch.Vertex2DSize)
}

// indexSliceToBytes reinterprets a []uint16 as a []byte without copying.
func indexSliceToBytes(indices []uint16) []byte {
	if len(indices) == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(&indices[0])), len(indices)*2)
}
