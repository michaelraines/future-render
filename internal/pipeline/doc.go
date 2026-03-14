// Package pipeline defines the render pipeline model: passes, draw commands,
// and their execution. Each pass has declared input and output resources and
// a deterministic execution function with no hidden state.
//
// Even the MVP 2D implementation uses this model so that adding 3D passes
// later is additive, not a rewrite.
package pipeline
