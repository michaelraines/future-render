package futurerender

// This file provides the public input API, wrapping the internal input system
// with user-friendly functions matching Ebitengine's API.

// IsKeyPressed returns whether the given key is currently pressed.
func IsKeyPressed(key Key) bool {
	if globalEngine == nil {
		return false
	}
	// Will delegate to internal input state once platform integration is done
	return false
}

// InputChars returns the Unicode characters input since the last frame.
func InputChars() []rune {
	return nil
}

// IsMouseButtonPressed returns whether the given mouse button is pressed.
func IsMouseButtonPressed(button MouseButton) bool {
	return false
}

// CursorPosition returns the current cursor position in logical pixels.
func CursorPosition() (int, int) {
	return 0, 0
}

// Wheel returns the mouse wheel delta since the last frame.
func Wheel() (float64, float64) {
	return 0, 0
}

// TouchIDs returns the IDs of all active touch points.
func TouchIDs() []TouchID {
	return nil
}

// TouchPosition returns the position of a touch point.
func TouchPosition(id TouchID) (int, int) {
	return 0, 0
}

// GamepadIDs returns the IDs of connected gamepads.
func GamepadIDs() []GamepadID {
	return nil
}

// GamepadAxisValue returns the value of a gamepad axis.
func GamepadAxisValue(id GamepadID, axis int) float64 {
	return 0
}

// IsGamepadButtonPressed returns whether a gamepad button is pressed.
func IsGamepadButtonPressed(id GamepadID, button GamepadButton) bool {
	return false
}

// Key represents a keyboard key. Maps to platform.Key values.
type Key int

// Key constants matching ebiten's key names.
const (
	KeyA Key = iota
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySpace
	KeyEnter
	KeyEscape
	KeyTab
	KeyBackspace
	KeyUp
	KeyDown
	KeyLeftKey
	KeyRightKey
	KeyShift
	KeyControl
	KeyAlt
)

// MouseButton represents a mouse button.
type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
	MouseButtonMiddle
)

// TouchID represents a touch point identifier.
type TouchID int

// GamepadID represents a gamepad identifier.
type GamepadID int

// GamepadButton represents a gamepad button.
type GamepadButton int
