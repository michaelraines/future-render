// Package math provides vector, matrix, quaternion, color, and geometric
// primitive types for the rendering engine. All types are value types with
// no heap allocations in their operations. This package is pure Go with no
// CGo dependencies.
//
// The coordinate system is right-handed with Y-up for 3D operations.
// 2D operations use screen coordinates (origin top-left, Y-down) unless
// otherwise documented.
//
// All angle parameters are in radians unless otherwise documented.
package math
