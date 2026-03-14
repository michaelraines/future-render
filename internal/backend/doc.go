// Package backend defines the graphics backend abstraction interface.
//
// This interface isolates all GPU-specific operations behind a set of
// narrow interfaces that can be implemented by Metal, Vulkan, OpenGL,
// WebGL/WebGPU, or Direct3D backends. No backend-specific types, constants,
// or calls should leak above this layer.
//
// The interface is designed with 3D forward-compatibility in mind: it supports
// depth buffers, 3D vertex formats, and perspective projection even though
// Phase 1 only uses 2D rendering.
package backend
