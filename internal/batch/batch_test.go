package batch

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/michaelraines/future-render/internal/backend"
)

func TestBatcherMerge(t *testing.T) {
	b := NewBatcher(65535, 65535)

	// Two quads with the same state should merge into one batch
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.AddQuad(20, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)

	batches := b.Flush()
	require.Len(t, batches, 1)
	require.Len(t, batches[0].Vertices, 8)
	require.Len(t, batches[0].Indices, 12)
}

func TestBatcherSplit(t *testing.T) {
	b := NewBatcher(65535, 65535)

	// Different textures should produce separate batches
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 2, backend.BlendSourceOver, 0)

	batches := b.Flush()
	require.Len(t, batches, 2)
}

func TestBatcherBlendModeSplit(t *testing.T) {
	b := NewBatcher(65535, 65535)

	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendAdditive, 0)

	batches := b.Flush()
	require.Len(t, batches, 2)
}

func TestBatcherFilterSplit(t *testing.T) {
	b := NewBatcher(65535, 65535)

	// Same texture but different filters should produce separate batches
	b.Add(DrawCommand{
		Vertices:  []Vertex2D{{PosX: 0, PosY: 0}, {PosX: 10, PosY: 0}, {PosX: 10, PosY: 10}, {PosX: 0, PosY: 10}},
		Indices:   []uint16{0, 1, 2, 0, 2, 3},
		TextureID: 1,
		BlendMode: backend.BlendSourceOver,
		Filter:    backend.FilterNearest,
	})
	b.Add(DrawCommand{
		Vertices:  []Vertex2D{{PosX: 20, PosY: 0}, {PosX: 30, PosY: 0}, {PosX: 30, PosY: 10}, {PosX: 20, PosY: 10}},
		Indices:   []uint16{0, 1, 2, 0, 2, 3},
		TextureID: 1,
		BlendMode: backend.BlendSourceOver,
		Filter:    backend.FilterLinear,
	})

	batches := b.Flush()
	require.Len(t, batches, 2)
	require.Equal(t, backend.FilterNearest, batches[0].Filter)
	require.Equal(t, backend.FilterLinear, batches[1].Filter)
}

func TestBatcherFilterMerge(t *testing.T) {
	b := NewBatcher(65535, 65535)

	// Same texture and same filter should merge
	b.Add(DrawCommand{
		Vertices:  []Vertex2D{{PosX: 0, PosY: 0}, {PosX: 10, PosY: 0}, {PosX: 10, PosY: 10}, {PosX: 0, PosY: 10}},
		Indices:   []uint16{0, 1, 2, 0, 2, 3},
		TextureID: 1,
		BlendMode: backend.BlendSourceOver,
		Filter:    backend.FilterLinear,
	})
	b.Add(DrawCommand{
		Vertices:  []Vertex2D{{PosX: 20, PosY: 0}, {PosX: 30, PosY: 0}, {PosX: 30, PosY: 10}, {PosX: 20, PosY: 10}},
		Indices:   []uint16{0, 1, 2, 0, 2, 3},
		TextureID: 1,
		BlendMode: backend.BlendSourceOver,
		Filter:    backend.FilterLinear,
	})

	batches := b.Flush()
	require.Len(t, batches, 1)
	require.Equal(t, backend.FilterLinear, batches[0].Filter)
	require.Len(t, batches[0].Vertices, 8)
}

func TestBatcherReset(t *testing.T) {
	b := NewBatcher(65535, 65535)
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.Reset()
	require.Equal(t, 0, b.CommandCount())
}

func TestBatcherIndexOffset(t *testing.T) {
	b := NewBatcher(65535, 65535)

	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.AddQuad(20, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)

	batches := b.Flush()
	indices := batches[0].Indices
	// Second quad indices should be offset by 4 (first quad's vertex count)
	require.Equal(t, uint16(4), indices[6])
	require.Equal(t, uint16(5), indices[7])
	require.Equal(t, uint16(6), indices[8])
}

func TestVertex2DFormat(t *testing.T) {
	f := Vertex2DFormat()
	require.Equal(t, Vertex2DSize, f.Stride)
	require.Len(t, f.Attributes, 3)

	require.Equal(t, "position", f.Attributes[0].Name)
	require.Equal(t, backend.AttributeFloat2, f.Attributes[0].Format)
	require.Equal(t, 0, f.Attributes[0].Offset)

	require.Equal(t, "texcoord", f.Attributes[1].Name)
	require.Equal(t, backend.AttributeFloat2, f.Attributes[1].Format)
	require.Equal(t, 8, f.Attributes[1].Offset)

	require.Equal(t, "color", f.Attributes[2].Name)
	require.Equal(t, backend.AttributeFloat4, f.Attributes[2].Format)
	require.Equal(t, 16, f.Attributes[2].Offset)
}

func TestBatcherFlushEmpty(t *testing.T) {
	b := NewBatcher(65535, 65535)
	batches := b.Flush()
	require.Equal(t, []Batch(nil), batches)
}

func TestBatcherSortOrder(t *testing.T) {
	b := NewBatcher(65535, 65535)

	// Add commands in reverse order: high shader/blend/texture first
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 3, backend.BlendAdditive, 2)
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 2, backend.BlendSourceOver, 0)
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 1)

	batches := b.Flush()

	// Should be sorted by shader, then blend, then texture
	// shader=0, blend=SourceOver(1), tex=1
	// shader=0, blend=SourceOver(1), tex=2
	// shader=1, blend=SourceOver(1), tex=1
	// shader=2, blend=Additive(2), tex=3
	require.Len(t, batches, 4)
	require.Equal(t, uint32(0), batches[0].ShaderID)
	require.Equal(t, backend.BlendSourceOver, batches[0].BlendMode)
	require.Equal(t, uint32(1), batches[0].TextureID)

	require.Equal(t, uint32(0), batches[1].ShaderID)
	require.Equal(t, backend.BlendSourceOver, batches[1].BlendMode)
	require.Equal(t, uint32(2), batches[1].TextureID)

	require.Equal(t, uint32(1), batches[2].ShaderID)
	require.Equal(t, backend.BlendSourceOver, batches[2].BlendMode)
	require.Equal(t, uint32(1), batches[2].TextureID)

	require.Equal(t, uint32(2), batches[3].ShaderID)
	require.Equal(t, backend.BlendAdditive, batches[3].BlendMode)
	require.Equal(t, uint32(3), batches[3].TextureID)
}

func BenchmarkBatcherFlush100Quads(b *testing.B) {
	batcher := NewBatcher(65535, 65535)
	for b.Loop() {
		for i := 0; i < 100; i++ {
			batcher.AddQuad(
				float32(i*12), 0, 10, 10,
				0, 0, 1, 1,
				1, 1, 1, 1,
				1, backend.BlendSourceOver, 0,
			)
		}
		_ = batcher.Flush()
	}
}

func BenchmarkBatcherFlush1000QuadsMixed(b *testing.B) {
	batcher := NewBatcher(65535, 65535)
	for b.Loop() {
		for i := 0; i < 1000; i++ {
			texID := uint32(i % 8)
			batcher.AddQuad(
				float32(i*12), 0, 10, 10,
				0, 0, 1, 1,
				1, 1, 1, 1,
				texID, backend.BlendSourceOver, 0,
			)
		}
		_ = batcher.Flush()
	}
}
