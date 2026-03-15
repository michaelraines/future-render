// Package softdelegate provides an embeddable Encoder and unwrapping
// interfaces for soft-delegating GPU backends. Each backend that wraps
// the software rasterizer embeds softdelegate.Encoder to eliminate
// boilerplate delegation code.
//
// When a backend is converted to real GPU bindings, the softdelegate
// embed is replaced by a direct implementation.
package softdelegate
