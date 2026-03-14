# Future Render — Agent Directives

## Project Overview

Future Render is a production-grade 2D/3D rendering engine in pure Go. Phase 1
targets full 2D feature parity with Ebitengine. The architecture is designed
from day one to support 3D rendering in later phases without rewrites.

Key documents:
- `DESIGN.md` — architecture, layer diagram, API design rationale
- `RESEARCH.md` — survey of Ebitengine, Raylib, bgfx, wgpu, Godot, Bevy, Three.js
- `FUTURE_3D.md` — 3D integration plan and Phase 1 constraints
- `ROADMAP.md` — phased implementation plan (update as work progresses)

## Build & Test

```bash
# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run benchmarks
go test -bench=. ./math/ ./internal/batch/

# Vet
go vet ./...

# Build (core library — no platform deps yet)
go build ./...
```

There are no external dependencies yet (`go.mod` has only the standard
library). When platform backends are added, they will use build tags
(`-tags glfw`).

## Architecture Rules

These are non-negotiable. Violating them creates technical debt that compounds.

1. **Layer direction is strictly downward.** No package may import from a layer
   above it. The layers top-to-bottom: public API (`engine.go`, `image.go`,
   `input.go`) → `internal/pipeline` → `internal/batch` →
   `internal/backend` → `internal/platform`.

2. **Backend types never leak to game code.** The public API uses
   engine-specific types (`BlendMode`, `Filter`) that map to internal backend
   types. Users never import `internal/`.

3. **No 2D-only assumptions in internal layers.** The backend, pipeline, and
   batch systems must work for both 2D and 3D. Read `FUTURE_3D.md` "What
   Phase 1 Must NOT Do" before changing internal packages.

4. **No CGo in core packages.** `math/`, `internal/batch/`,
   `internal/pipeline/`, `internal/input/` must remain pure Go. CGo is
   permitted only in `internal/backend/<impl>/` and `internal/platform/<impl>/`.

5. **Interfaces are defined by consumers, not implementors.** Follow Go
   interface design conventions. Keep interfaces small and focused.

## Development Workflow

Follow this cycle for every change:

### 1. Understand Before Changing
- Read the relevant source files before modifying them
- Check `DESIGN.md` to understand where the change fits architecturally
- Check `FUTURE_3D.md` constraints if touching internal packages

### 2. Implement
- Make the minimal change needed
- Prefer editing existing files over creating new ones
- Don't add features, abstractions, or "improvements" beyond what was asked
- No empty files, placeholder packages, or premature abstractions

### 3. Test
- Run `go test ./...` after every change
- Run `go vet ./...` to catch issues
- Add tests for new logic — particularly in `math/`, `internal/batch/`,
  and `internal/pipeline/`
- Tests must pass before committing

### 4. Verify Build
- Run `go build ./...` to ensure everything compiles
- If adding platform-specific code, verify build tags work

### 5. Update Docs
- Update `ROADMAP.md` when completing milestone tasks
- Don't create new markdown files unless explicitly asked

### Loop: test → fix → test
If tests fail, fix the issue and re-run. Don't commit broken code. Don't skip
tests. Don't use `-count=0` or other tricks to hide failures.

## Code Style

- Standard Go formatting (`gofmt`)
- Error returns use `(T, error)` pattern, not panics
- Exported types and functions have doc comments
- Internal packages use `internal/` path convention
- Test files are `*_test.go` in the same package
- Benchmarks use `Benchmark*` naming convention

## Naming Conventions

- Public API types match Ebitengine where applicable: `Game`, `Image`,
  `GeoM`, `DrawImageOptions`, `Vertex`, `Key`, `MouseButton`
- Math types use short names: `Vec2`, `Vec3`, `Mat3`, `Mat4`, `Quat`
- Backend interfaces use GPU terminology: `Device`, `Texture`, `Buffer`,
  `Shader`, `Pipeline`, `CommandEncoder`
- Platform interfaces: `Window`, `InputHandler`

## Commit Messages

- Use imperative mood: "Add sprite pass" not "Added sprite pass"
- First line under 72 characters
- Reference the milestone when relevant: "M2: wire DrawImage to batcher"

## Common Pitfalls

- **Don't hardcode orthographic projection** in pipeline internals — projection
  matrix must be a parameter
- **Don't assume Vertex2D is the only format** — batcher and pipeline must
  support arbitrary vertex formats
- **Don't tie render targets to screen size** — off-screen targets can be any
  dimension
- **Don't remove depth/3D fields** from backend types even though Phase 1
  doesn't use them
- **Don't merge pipeline and backend layers** — their separation is essential
  for 3D
- **Don't add Ebitengine as a dependency** — this is a clean-room implementation

## Test Coverage Expectations

| Package | Coverage Goal | Notes |
|---|---|---|
| `math/` | High | Pure functions, easy to test exhaustively |
| `internal/batch/` | High | Core optimization logic, must be correct |
| `internal/pipeline/` | Medium | Test pass ordering and context propagation |
| `internal/input/` | Medium | Test state transitions, edge detection |
| `internal/backend/` | Low (interfaces) | Implementations tested via integration |
| `internal/platform/` | Low (interfaces) | Implementations tested via integration |
| Public API | Medium | Test option defaults, GeoM transforms |
