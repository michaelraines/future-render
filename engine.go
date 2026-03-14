// Package futurerender is a production-grade 2D/3D rendering engine for Go.
//
// The engine provides an API compatible with Ebitengine's game loop model:
// a Game interface with Update(), Draw(), and Layout() methods. The engine
// manages the window, input, audio, and rendering pipeline.
//
// Basic usage:
//
//	type MyGame struct{}
//
//	func (g *MyGame) Update() error { return nil }
//	func (g *MyGame) Draw(screen *futurerender.Image) {}
//	func (g *MyGame) Layout(outsideWidth, outsideHeight int) (int, int) {
//	    return 320, 240
//	}
//
//	func main() {
//	    futurerender.RunGame(&MyGame{})
//	}
package futurerender

import (
	"errors"
	"sync/atomic"
	"time"
)

// Game is the interface that game implementations must satisfy.
// This matches Ebitengine's Game interface for compatibility.
type Game interface {
	// Update is called every tick. Game logic goes here.
	// Return a non-nil error to terminate the game loop.
	Update() error

	// Draw is called every frame. Rendering goes here.
	// The screen image is the render target for the frame.
	Draw(screen *Image)

	// Layout accepts the outside (window) size and returns the logical
	// screen size. The engine scales the logical screen to fit the window.
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}

// Termination is returned from Update() to cleanly exit the game loop.
var Termination = errors.New("game terminated")

// RunGame starts the game loop with the given Game implementation.
// This function blocks until the game exits. It must be called from
// the main goroutine on platforms that require it (macOS, iOS).
func RunGame(game Game) error {
	e := newEngine(game)
	return e.run()
}

// SetWindowSize sets the window size in logical pixels.
func SetWindowSize(width, height int) {
	if globalEngine != nil {
		globalEngine.setWindowSize(width, height)
	}
}

// SetWindowTitle sets the window title.
func SetWindowTitle(title string) {
	if globalEngine != nil {
		globalEngine.setWindowTitle(title)
	}
}

// SetFullscreen sets fullscreen mode.
func SetFullscreen(fullscreen bool) {
	if globalEngine != nil {
		globalEngine.setFullscreen(fullscreen)
	}
}

// IsFullscreen returns whether the window is in fullscreen mode.
func IsFullscreen() bool {
	if globalEngine != nil {
		return globalEngine.isFullscreen()
	}
	return false
}

// SetMaxTPS sets the maximum ticks per second. The default is 60.
// Set to 0 for uncapped TPS (sync to frame rate).
func SetMaxTPS(tps int) {
	if tps < 0 {
		tps = 0
	}
	maxTPS.Store(int64(tps))
}

// MaxTPS returns the current maximum ticks per second.
func MaxTPS() int {
	return int(maxTPS.Load())
}

// SetVsyncEnabled enables or disables vertical synchronization.
func SetVsyncEnabled(enabled bool) {
	if globalEngine != nil {
		globalEngine.setVSync(enabled)
	}
}

// IsVsyncEnabled returns whether VSync is enabled.
func IsVsyncEnabled() bool {
	if globalEngine != nil {
		return globalEngine.isVSync()
	}
	return true
}

// CurrentFPS returns the current frames per second.
func CurrentFPS() float64 {
	if globalEngine != nil {
		return globalEngine.currentFPS()
	}
	return 0
}

// CurrentTPS returns the current ticks per second.
func CurrentTPS() float64 {
	if globalEngine != nil {
		return globalEngine.currentTPS()
	}
	return 0
}

// SetCursorMode sets the cursor visibility and lock mode.
func SetCursorMode(mode CursorMode) {
	if globalEngine != nil {
		globalEngine.setCursorMode(mode)
	}
}

// CursorMode constants.
type CursorMode int

const (
	CursorModeVisible  CursorMode = iota // Normal cursor
	CursorModeHidden                     // Hidden cursor
	CursorModeCaptured                   // Hidden and locked to window
)

// DeviceScaleFactor returns the device pixel ratio.
func DeviceScaleFactor() float64 {
	if globalEngine != nil {
		return globalEngine.deviceScaleFactor()
	}
	return 1.0
}

// --- Engine internals ---

var (
	globalEngine *engine
	maxTPS       atomic.Int64
)

func init() {
	maxTPS.Store(60)
}

type engine struct {
	game Game

	// Timing
	lastUpdateTime time.Time
	lastDrawTime   time.Time
	updateCount    int64
	drawCount      int64
	fpsAccum       float64
	tpsAccum       float64
	fpsValue       float64
	tpsValue       float64
	fpsTimer       time.Time
}

func newEngine(game Game) *engine {
	e := &engine{
		game: game,
	}
	globalEngine = e
	return e
}

func (e *engine) run() error {
	// The actual run implementation will integrate with the platform window
	// and backend. For now, this provides the timing/loop structure.
	//
	// The real implementation will:
	// 1. Create a platform window
	// 2. Initialize the graphics backend
	// 3. Set up the render pipeline
	// 4. Run the main loop with fixed-timestep update + variable draw
	//
	// This placeholder ensures the public API compiles and the package
	// structure is valid.
	return errors.New("engine: no platform backend available (build with platform tags)")
}

func (e *engine) setWindowSize(width, height int)  {}
func (e *engine) setWindowTitle(title string)       {}
func (e *engine) setFullscreen(fullscreen bool)     {}
func (e *engine) isFullscreen() bool                { return false }
func (e *engine) setVSync(enabled bool)             {}
func (e *engine) isVSync() bool                     { return true }
func (e *engine) currentFPS() float64               { return e.fpsValue }
func (e *engine) currentTPS() float64               { return e.tpsValue }
func (e *engine) setCursorMode(mode CursorMode)     {}
func (e *engine) deviceScaleFactor() float64        { return 1.0 }
