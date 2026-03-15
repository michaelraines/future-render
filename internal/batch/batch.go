package batch

import (
	"math"
	"sort"

	"github.com/michaelraines/future-render/internal/backend"
)

// Vertex2D represents a single vertex in a 2D draw call.
// This is the standard vertex format for Phase 1 2D rendering.
type Vertex2D struct {
	PosX, PosY float32 // position
	TexU, TexV float32 // texture coordinates
	R, G, B, A float32 // vertex color
}

// Vertex2DSize is the byte size of a Vertex2D.
const Vertex2DSize = 32 // 8 float32s × 4 bytes

// Vertex2DFormat returns the VertexFormat for Vertex2D.
func Vertex2DFormat() backend.VertexFormat {
	return backend.VertexFormat{
		Stride: Vertex2DSize,
		Attributes: []backend.VertexAttribute{
			{Name: "position", Format: backend.AttributeFloat2, Offset: 0},
			{Name: "texcoord", Format: backend.AttributeFloat2, Offset: 8},
			{Name: "color", Format: backend.AttributeFloat4, Offset: 16},
		},
	}
}

// DrawCommand represents a single draw command before batching.
type DrawCommand struct {
	Vertices  []Vertex2D
	Indices   []uint16
	TextureID uint32 // opaque texture identifier for sorting
	BlendMode backend.BlendMode
	Filter    backend.TextureFilter // texture filter (nearest or linear)
	ShaderID  uint32                // opaque shader identifier for sorting
	Depth     float32               // sort key for back-to-front or front-to-back ordering
}

// Batch represents a group of draw commands that share the same state.
type Batch struct {
	Vertices  []Vertex2D
	Indices   []uint16
	TextureID uint32
	BlendMode backend.BlendMode
	Filter    backend.TextureFilter
	ShaderID  uint32
}

// Batcher accumulates draw commands and produces optimized batches.
type Batcher struct {
	commands    []DrawCommand
	maxVertices int
	maxIndices  int
}

// NewBatcher creates a new Batcher with the given capacity hints.
func NewBatcher(maxVertices, maxIndices int) *Batcher {
	return &Batcher{
		commands:    make([]DrawCommand, 0, 256),
		maxVertices: maxVertices,
		maxIndices:  maxIndices,
	}
}

// Add adds a draw command to be batched.
func (b *Batcher) Add(cmd DrawCommand) {
	b.commands = append(b.commands, cmd)
}

// AddQuad is a convenience method that adds a textured quad.
func (b *Batcher) AddQuad(
	x, y, w, h float32,
	u0, v0, u1, v1 float32,
	r, g, bl, a float32,
	textureID uint32,
	blendMode backend.BlendMode,
	shaderID uint32,
) {
	baseIdx := uint16(0) // will be adjusted during batching
	b.commands = append(b.commands, DrawCommand{
		Vertices: []Vertex2D{
			{PosX: x, PosY: y, TexU: u0, TexV: v0, R: r, G: g, B: bl, A: a},
			{PosX: x + w, PosY: y, TexU: u1, TexV: v0, R: r, G: g, B: bl, A: a},
			{PosX: x + w, PosY: y + h, TexU: u1, TexV: v1, R: r, G: g, B: bl, A: a},
			{PosX: x, PosY: y + h, TexU: u0, TexV: v1, R: r, G: g, B: bl, A: a},
		},
		Indices:   []uint16{baseIdx, baseIdx + 1, baseIdx + 2, baseIdx, baseIdx + 2, baseIdx + 3},
		TextureID: textureID,
		BlendMode: blendMode,
		ShaderID:  shaderID,
	})
}

// Flush produces batches from accumulated commands and resets the batcher.
// Commands are sorted by (shader, blend mode, texture) to minimize state changes,
// then merged when adjacent commands share the same state.
func (b *Batcher) Flush() []Batch {
	if len(b.commands) == 0 {
		return nil
	}

	// Sort for optimal batching: group by state
	sort.Slice(b.commands, func(i, j int) bool {
		ci, cj := b.commands[i], b.commands[j]
		if ci.ShaderID != cj.ShaderID {
			return ci.ShaderID < cj.ShaderID
		}
		if ci.BlendMode != cj.BlendMode {
			return ci.BlendMode < cj.BlendMode
		}
		if ci.Filter != cj.Filter {
			return ci.Filter < cj.Filter
		}
		if ci.TextureID != cj.TextureID {
			return ci.TextureID < cj.TextureID
		}
		return ci.Depth < cj.Depth
	})

	batches := make([]Batch, 0, 16)
	var current *Batch

	for i := range b.commands {
		cmd := &b.commands[i]

		// Check if we can merge with the current batch
		canMerge := current != nil &&
			current.TextureID == cmd.TextureID &&
			current.BlendMode == cmd.BlendMode &&
			current.Filter == cmd.Filter &&
			current.ShaderID == cmd.ShaderID &&
			len(current.Vertices)+len(cmd.Vertices) <= b.maxVertices &&
			len(current.Indices)+len(cmd.Indices) <= b.maxIndices &&
			len(current.Vertices)+len(cmd.Vertices) <= math.MaxUint16

		if canMerge {
			// Merge: adjust indices and append
			vertexOffset := uint16(len(current.Vertices))
			for _, idx := range cmd.Indices {
				current.Indices = append(current.Indices, idx+vertexOffset)
			}
			current.Vertices = append(current.Vertices, cmd.Vertices...)
		} else {
			// Start a new batch
			batches = append(batches, Batch{
				Vertices:  make([]Vertex2D, len(cmd.Vertices)),
				Indices:   make([]uint16, len(cmd.Indices)),
				TextureID: cmd.TextureID,
				BlendMode: cmd.BlendMode,
				Filter:    cmd.Filter,
				ShaderID:  cmd.ShaderID,
			})
			current = &batches[len(batches)-1]
			copy(current.Vertices, cmd.Vertices)
			copy(current.Indices, cmd.Indices)
		}
	}

	// Reset for next frame
	b.commands = b.commands[:0]

	return batches
}

// Reset clears accumulated commands without producing batches.
func (b *Batcher) Reset() {
	b.commands = b.commands[:0]
}

// CommandCount returns the number of pending commands.
func (b *Batcher) CommandCount() int {
	return len(b.commands)
}
