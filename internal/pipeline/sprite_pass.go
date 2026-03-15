package pipeline

import (
	"unsafe"

	"github.com/michaelraines/future-render/internal/backend"
	"github.com/michaelraines/future-render/internal/batch"
)

// TextureResolver maps a batcher texture ID to a backend.Texture.
type TextureResolver func(textureID uint32) backend.Texture

// ShaderInfo holds a custom shader's backend resources for rendering.
type ShaderInfo struct {
	Shader   backend.Shader
	Pipeline backend.Pipeline
}

// ShaderResolver maps a batcher shader ID to a ShaderInfo.
// Returns nil if the shader ID is not registered (use default).
type ShaderResolver func(shaderID uint32) *ShaderInfo

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

	// ResolveShader maps batch shader IDs to custom shader info.
	ResolveShader ShaderResolver

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

	// Track which shader is currently bound to minimize state changes.
	currentShaderID := uint32(0)

	// Set default pipeline and projection uniform.
	enc.SetPipeline(sp.pipeline)
	sp.shader.SetUniformMat4("uProjection", sp.Projection)
	sp.shader.SetUniformInt("uTexture", 0)

	for i := range batches {
		b := &batches[i]

		// Switch shader if this batch uses a different one.
		if b.ShaderID != currentShaderID {
			if b.ShaderID == 0 {
				// Switch back to default shader.
				enc.SetPipeline(sp.pipeline)
				sp.shader.SetUniformMat4("uProjection", sp.Projection)
				sp.shader.SetUniformInt("uTexture", 0)
			} else if sp.ResolveShader != nil {
				info := sp.ResolveShader(b.ShaderID)
				if info != nil {
					enc.SetPipeline(info.Pipeline)
					info.Shader.SetUniformMat4("uProjection", sp.Projection)
				}
			}
			currentShaderID = b.ShaderID
		}

		// Upload vertex data.
		vertexData := vertexSliceToBytes(b.Vertices)
		sp.vertexBuf.Upload(vertexData)
		enc.SetVertexBuffer(sp.vertexBuf, 0)

		// Upload index data.
		indexData := indexSliceToBytes(b.Indices)
		sp.indexBuf.Upload(indexData)
		enc.SetIndexBuffer(sp.indexBuf, backend.IndexUint16)

		// Bind texture and set per-draw filter via sampler object.
		if sp.ResolveTexture != nil {
			tex := sp.ResolveTexture(b.TextureID)
			if tex != nil {
				enc.SetTexture(tex, 0)
			}
		}
		enc.SetTextureFilter(0, b.Filter)

		// Handle fill rule.
		if b.FillRule == backend.FillRuleEvenOdd {
			sp.drawEvenOdd(enc, b)
		} else {
			enc.DrawIndexed(len(b.Indices), 1, 0)
		}
	}
}

// drawEvenOdd renders a batch using the even-odd fill rule via stencil.
// Pass 1: draw triangles to stencil only (INVERT), color writes disabled.
// Pass 2: redraw with stencil test NOTEQUAL 0, then reset stencil state.
func (sp *SpritePass) drawEvenOdd(enc backend.CommandEncoder, b *batch.Batch) {
	// Pass 1: write to stencil, no color output.
	enc.SetColorWrite(false)
	enc.SetStencil(true, backend.StencilDescriptor{
		Func:      backend.CompareAlways,
		Ref:       0,
		Mask:      0xFF,
		SFail:     backend.StencilKeep,
		DPFail:    backend.StencilKeep,
		DPPass:    backend.StencilInvert,
		WriteMask: 0xFF,
	})
	enc.DrawIndexed(len(b.Indices), 1, 0)

	// Pass 2: draw where stencil != 0.
	enc.SetColorWrite(true)
	enc.SetStencil(true, backend.StencilDescriptor{
		Func:      backend.CompareNotEqual,
		Ref:       0,
		Mask:      0xFF,
		SFail:     backend.StencilKeep,
		DPFail:    backend.StencilKeep,
		DPPass:    backend.StencilZero, // clear stencil as we draw
		WriteMask: 0xFF,
	})
	enc.DrawIndexed(len(b.Indices), 1, 0)

	// Disable stencil for subsequent batches.
	enc.SetStencil(false, backend.StencilDescriptor{})
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
