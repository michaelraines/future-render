package batch

import (
	"testing"

	"github.com/michaelraines/future-render/internal/backend"
)

func TestBatcherMerge(t *testing.T) {
	b := NewBatcher(65535, 65535)

	// Two quads with the same state should merge into one batch
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.AddQuad(20, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)

	batches := b.Flush()
	if len(batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}
	if len(batches[0].Vertices) != 8 {
		t.Errorf("expected 8 vertices, got %d", len(batches[0].Vertices))
	}
	if len(batches[0].Indices) != 12 {
		t.Errorf("expected 12 indices, got %d", len(batches[0].Indices))
	}
}

func TestBatcherSplit(t *testing.T) {
	b := NewBatcher(65535, 65535)

	// Different textures should produce separate batches
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 2, backend.BlendSourceOver, 0)

	batches := b.Flush()
	if len(batches) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(batches))
	}
}

func TestBatcherBlendModeSplit(t *testing.T) {
	b := NewBatcher(65535, 65535)

	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendAdditive, 0)

	batches := b.Flush()
	if len(batches) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(batches))
	}
}

func TestBatcherReset(t *testing.T) {
	b := NewBatcher(65535, 65535)
	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.Reset()
	if b.CommandCount() != 0 {
		t.Errorf("expected 0 commands after reset, got %d", b.CommandCount())
	}
}

func TestBatcherIndexOffset(t *testing.T) {
	b := NewBatcher(65535, 65535)

	b.AddQuad(0, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)
	b.AddQuad(20, 0, 10, 10, 0, 0, 1, 1, 1, 1, 1, 1, 1, backend.BlendSourceOver, 0)

	batches := b.Flush()
	indices := batches[0].Indices
	// Second quad indices should be offset by 4 (first quad's vertex count)
	if indices[6] != 4 || indices[7] != 5 || indices[8] != 6 {
		t.Errorf("second quad indices not properly offset: %v", indices[6:12])
	}
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
