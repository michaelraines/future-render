# Rendering Engine Research

Technical survey of existing rendering engines to inform the design of Future Render,
a Go-based 2D/3D rendering engine.

---

## 1. Ebitengine (Go) — Primary API Reference

### Architecture

Ebitengine is the most mature 2D game engine in Go. Its internal rendering stack is deeply layered:

1. **`buffered`** — Caches CPU-side pixel data, syncs to GPU lazily.
2. **`atlas`** — Automatic texture atlas management (4096x4096 default). Images are packed
   with 1px padding. Sub-images share atlas allocations.
3. **`restorable`** — GPU context loss recovery. Records all draw operations for replay
   if the context is lost (essential for mobile/WebGL).
4. **`graphicscommand`** — Queues draw commands without immediate execution.
   `DrawTriangles()` enqueues via `theCommandQueueManager`. Commands are flushed in batch.
5. **`graphicsdriver`** — Backend abstraction interface. Implementations: OpenGL, Metal,
   DirectX 11/12, PlayStation 5.

### Game Loop

The `Game` interface defines three methods:
- `Update()` — Fixed timestep logic tick (default 60 TPS, configurable).
- `Draw(screen *ebiten.Image)` — Variable-rate rendering.
- `Layout(outsideWidth, outsideHeight int) (int, int)` — Logical-to-physical size mapping.

TPS is prioritized over FPS: under lag, Draw is skipped to maintain simulation speed.
Multiple Updates can fire before a single Draw.

### Rendering Model

- Everything is an `ebiten.Image`. Drawing = compositing one image onto another.
- `DrawImage` internally generates 2 triangles (a quad).
- Successive draw calls are automatically merged when source images share the same atlas
  and no state changes occur. The sprites example draws 10,000+ sprites in 1-2 GPU draw calls.
- Batching breaks on: atlas boundary crossings, blend mode changes, shader changes,
  offscreen targets.

### Shader System (Kage)

Kage is a Go-syntax-compatible shading language for fragment shaders only (no vertex shaders).
Compiled at runtime and transpiled to GLSL/HLSL/MSL depending on the active backend.
Supports up to 4 source textures per shader invocation. No struct types.

### Key Design Tradeoffs

- **Simplicity over flexibility**: No scene graph, no ECS, no vertex shaders, no compute.
- **Automatic atlas**: Great for simple games, limits fine-grained control.
- **Context loss recovery**: Adds memory/complexity overhead but essential for mobile.
- **2D only**: Community 3D projects (Tetra3D) are limited by the 2D architecture.
- **Thread-safe by default**: Overhead removable via `ebitenginesinglethread` build tag.

---

## 2. Raylib (C) — Simplicity Reference

### 2D/3D Unification

Raylib unifies 2D and 3D through **rlgl**, a pseudo-OpenGL 1.1 immediate-mode abstraction.
All 2D shape drawing uses only 6 rlgl functions: `rlBegin()`, `rlEnd()`, `rlVertex3f()`,
`rlTexCoord2f()`, `rlNormal3f()`, `rlSetTexture()`. 2D drawing is simply 3D drawing with
an orthographic projection — there is no separate 2D pipeline.

### API Design Philosophy

- Procedural C (C99), no OOP. All structures are transparent (direct field access).
- Immediate-mode rendering: `BeginDrawing()` -> draw calls -> `EndDrawing()`.
- Modular single-header libraries: rlgl, raymath, rcamera, rgestures are standalone.
- Intentionally limited scope: simplicity and learnability over capability.

### Key Tradeoffs

- **OpenGL only**: No Vulkan, Metal, or D3D. rlgl's immediate-mode model fundamentally
  conflicts with Vulkan's explicit command-buffer approach.
- **Single-threaded**: No multi-threading.
- **No draw call sorting**: Simple batching on state change only.

---

## 3. bgfx (C++) — Backend Abstraction Reference

### Multi-API Abstraction

bgfx abstracts over 9 backends: Noop, D3D 11/12, Metal, OpenGL/GLES, Vulkan, GNM (PS4),
AGC (PS5), NVN (Switch). All resources use opaque handle-based references. Backend selection
happens at init time.

All backends implement `RendererContextI`. The critical method is `submit(Frame*)`, which
dispatches all draw calls for the frame.

### Sort-Based Draw Call Bucketing

bgfx uses a **64-bit sort key** for draw call ordering:

- **Views** are the primary sort criterion (highest bits). Views = render passes/buckets.
- Within a view, sort mode is configurable:
  - `Default` — sort by shader/state (minimize state changes)
  - `DepthAscending/DepthDescending` — sort by depth
  - `Sequential` — preserve submission order (for GUI)
- Keys are sorted using **radix sort** before backend submission.

This design allows scene traversal once, with draw calls submitted for all passes
simultaneously. bgfx handles optimal ordering internally.

### Rendering Submission Model

- Command buffer architecture: API thread encodes commands, render thread processes.
- Up to 8 simultaneous encoding threads.
- Maximum 64K draw calls per frame.
- State is declarative: `setState()`, `setTexture()`, `setVertexBuffer()`, then `submit()`.

---

## 4. wgpu / WebGPU — Modern API Model Reference

### Core Concepts

WebGPU converges modern GPU API design (Vulkan, Metal, D3D12) into a portable abstraction:

- **Explicit resource management** (but less low-level than Vulkan)
- **Pipeline state objects** pre-compile all render state
- **Bind groups** for batched resource binding
- **Command buffer recording** separated from execution
- **First-class compute** alongside graphics

### Command Flow

`Device` -> `CommandEncoder` -> `RenderPass` / `ComputePass` -> `CommandBuffer` -> `Queue.submit()`

- **RenderPass**: Created from descriptor specifying color/depth attachments, load/store ops.
- **Pipeline State Object**: Bundles shaders, vertex layout, topology, multisampling, depth.
- **Bind Groups**: Resources bound in collections by update frequency (per-frame, per-material, per-object).

### wgpu Architecture (Three Layers)

1. **wgpu** — Safe API
2. **wgpu_core** — Validation and state tracking
3. **wgpu_hal** — Unsafe, zero-overhead hardware abstraction

Backend trait uses **static dispatch** for zero overhead. Naga handles shader cross-compilation
(WGSL <-> SPIR-V <-> HLSL <-> MSL <-> GLSL). A **pure-Go port of Naga** exists (github.com/gogpu/naga).

---

## 5. Godot (C++) — Render Server Reference

### RenderingServer Pattern

Godot's visual system uses the RenderingServer (RS) singleton. All GPU resources are
opaque RID handles. The scene tree never holds GPU state directly.

Thread safety via command-queue: calls from any thread are serialized into a lock-free queue;
the render thread drains once per frame.

### Layered Architecture

| Layer | Responsibility |
|---|---|
| Scene tree (Nodes) | Spatial hierarchy, visibility, game logic |
| RenderingServer | Resource management, draw-list building, culling |
| RendererCompositor | Backend selection, frame orchestration |
| RenderingDevice (RD) | Thin abstraction over Vulkan / D3D12 / Metal |

### Render Pipeline

Forward+ with clustered light culling (default desktop). Frustum divided into 3D grid of
clusters; each fragment looks up its cluster for relevant lights. O(fragments × lights_per_cluster)
vs O(fragments × total_lights).

### 2D + 3D Unification

Via Viewport: owns both a World3D scenario and Canvas layers. 3D rendered first, 2D composited
on top. Both pipelines feed into one viewport framebuffer through the same command queue.

---

## 6. Bevy (Rust) — Pipeline Architecture Reference

### Dual-World ECS

Main World (game logic) and Render World (GPU work) are separate ECS instances.
Pipelined rendering: frame N+1's game logic runs concurrently with frame N's GPU submission.

### Render Schedule

1. **Extract** — Sole sync point. Minimal data copied from Main to Render World.
2. **Prepare** — Upload to GPU buffers. Create bind groups.
3. **Queue** — Build render phases (Opaque3d, Transparent3d, etc.) as PhaseItems.
4. **PhaseSort** — Opaque front-to-back, transparent back-to-front.
5. **Render** — Execute render graph.
6. **Cleanup** — Despawn render world entities.

### Batching

- **Binned Phase Items**: Keyed by (pipeline, draw_function, material). Merged into single draws.
- **Sorted Phase Items**: Individually sorted for transparency.

---

## 7. Three.js (JavaScript) — API Design Reference

### Progressive Disclosure

Minimum viable program: ~5 lines. Scales to complex multi-pass rendering.
Deliberately unopinionated about app architecture.

### Abstraction Mapping

| GPU Concept | Three.js Abstraction |
|---|---|
| Vertex/index buffers | `BufferGeometry` |
| Shaders + uniforms | `Material` subclasses |
| Draw calls | `Mesh` = Geometry + Material |
| Transforms | `Object3D.position/.rotation/.scale` |

### Material Hierarchy

`MeshBasicMaterial` (unlit) -> `MeshLambertMaterial` (diffuse) -> `MeshPhongMaterial` (specular)
-> `MeshStandardMaterial` (PBR) -> `MeshPhysicalMaterial` (extended PBR) -> `ShaderMaterial` (custom).

---

## 8. Cross-Engine Comparison

| Dimension | Ebitengine | Raylib | bgfx | wgpu | Godot | Bevy | Three.js |
|---|---|---|---|---|---|---|---|
| **API style** | Game interface | Immediate-mode | Declarative submit | Command buffers | RID handles | ECS systems | Scene graph |
| **Backend** | GL/Metal/D3D | OpenGL only | 9 backends | Vk/Metal/D3D12/GL | Vk/D3D12/Metal/GLES3 | wgpu | WebGL/WebGPU |
| **Batching** | Auto (atlas) | On state change | 64-bit sort key | Manual | Auto (2D) | Binned phases | Manual |
| **Threading** | Safe by default | Single | Multi-encoder | Command buffers | Command queue | Pipelined worlds | Single |
| **2D + 3D** | 2D only | Unified (ortho) | Both | Both | Viewport | Separate phases | Both |
| **Shader** | Kage (frag only) | Raw GLSL | shaderc | Naga (WGSL) | Custom GLSL subset | WGSL | GLSL/WGSL |

---

## 9. Approach for Future Render

Based on this survey, Future Render adopts:

1. **Ebitengine-compatible Game interface** for the public API (Update/Draw/Layout).
2. **bgfx-inspired sort-based draw call batching** with a sort key encoding pipeline/texture/depth.
3. **wgpu/Godot-inspired backend abstraction** with Device, CommandEncoder, RenderPass, Pipeline
   interfaces — designed for modern explicit APIs but implementable on OpenGL.
4. **Bevy-inspired render pipeline model** with explicit passes that have declared inputs/outputs.
5. **Godot-inspired 2D/3D unification** where 2D is a specialized case of the 3D pipeline
   (orthographic projection, 2D vertex format).
6. **Three.js-inspired progressive API** that starts simple but allows dropping to lower levels.
7. **Pure Go math** (no CGo) with Vec2/Vec3/Vec4/Mat3/Mat4/Quat from day one.

What we deliberately avoid:
- Kage's fragment-only shader limitation (we will support vertex + fragment)
- Ebitengine's automatic atlas (user controls texture management)
- Raylib's OpenGL-only restriction (backend interface from day one)
- Three.js's lack of automatic batching (we batch by default)
- Bevy's ECS complexity (Go's interfaces + composition instead)
