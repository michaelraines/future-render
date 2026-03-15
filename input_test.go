package futurerender

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsKeyPressedNilEngine(t *testing.T) {
	withNilEngine(t)
	require.False(t, IsKeyPressed(KeyA))
	require.False(t, IsKeyPressed(KeySpace))
	require.False(t, IsKeyPressed(KeyEnter))
}

func TestInputCharsReturnsNil(t *testing.T) {
	require.Nil(t, InputChars())
}

func TestIsMouseButtonPressedReturnsFalse(t *testing.T) {
	require.False(t, IsMouseButtonPressed(MouseButtonLeft))
	require.False(t, IsMouseButtonPressed(MouseButtonRight))
	require.False(t, IsMouseButtonPressed(MouseButtonMiddle))
}

func TestCursorPositionReturnsZero(t *testing.T) {
	x, y := CursorPosition()
	require.Equal(t, 0, x)
	require.Equal(t, 0, y)
}

func TestWheelReturnsZero(t *testing.T) {
	xoff, yoff := Wheel()
	require.InDelta(t, 0.0, xoff, 1e-6)
	require.InDelta(t, 0.0, yoff, 1e-6)
}

func TestTouchIDsReturnsNil(t *testing.T) {
	require.Nil(t, TouchIDs())
}

func TestTouchPositionReturnsZero(t *testing.T) {
	x, y := TouchPosition(TouchID(0))
	require.Equal(t, 0, x)
	require.Equal(t, 0, y)
}

func TestGamepadIDsReturnsNil(t *testing.T) {
	require.Nil(t, GamepadIDs())
}

func TestGamepadAxisValueReturnsZero(t *testing.T) {
	val := GamepadAxisValue(GamepadID(0), 0)
	require.InDelta(t, 0.0, val, 1e-6)
}

func TestIsGamepadButtonPressedReturnsFalse(t *testing.T) {
	require.False(t, IsGamepadButtonPressed(GamepadID(0), GamepadButton(0)))
}

func TestKeyConstants(t *testing.T) {
	// Verify key constants are defined and distinct.
	require.NotEqual(t, KeyA, KeyB)
	require.NotEqual(t, KeySpace, KeyEnter)
}

func TestMouseButtonConstants(t *testing.T) {
	require.Equal(t, MouseButton(0), MouseButtonLeft)
	require.Equal(t, MouseButton(1), MouseButtonRight)
	require.Equal(t, MouseButton(2), MouseButtonMiddle)
}
