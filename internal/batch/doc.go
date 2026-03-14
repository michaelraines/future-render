// Package batch provides draw call batching and sorting for efficient
// GPU submission. It groups compatible draw calls (same texture, blend mode,
// shader) into batches to minimize state changes.
package batch
