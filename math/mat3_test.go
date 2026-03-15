package math

import (
	gomath "math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMat3Identity(t *testing.T) {
	m := Mat3Identity()
	v := m.MulVec2(NewVec2(3, 4))
	require.InDelta(t, 3.0, v.X, 1e-9)
	require.InDelta(t, 4.0, v.Y, 1e-9)
}

func TestMat3Translate(t *testing.T) {
	m := Mat3Translate(10, 20)
	v := m.MulVec2(NewVec2(1, 2))
	require.InDelta(t, 11.0, v.X, 1e-9)
	require.InDelta(t, 22.0, v.Y, 1e-9)
}

func TestMat3Scale(t *testing.T) {
	m := Mat3Scale(2, 3)
	v := m.MulVec2(NewVec2(5, 4))
	require.InDelta(t, 10.0, v.X, 1e-9)
	require.InDelta(t, 12.0, v.Y, 1e-9)
}

func TestMat3Rotate(t *testing.T) {
	m := Mat3Rotate(gomath.Pi / 2)
	v := m.MulVec2(NewVec2(1, 0))
	require.InDelta(t, 0.0, v.X, 1e-9)
	require.InDelta(t, 1.0, v.Y, 1e-9)
}

func TestMat3Shear(t *testing.T) {
	m := Mat3Shear(1, 0)
	v := m.MulVec2(NewVec2(0, 1))
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 1.0, v.Y, 1e-9)
}

func TestMat3Mul(t *testing.T) {
	a := Mat3Translate(5, 10)
	b := Mat3Identity()
	c := a.Mul(b)
	require.True(t, c.ApproxEqual(a, 1e-9))

	// Scale then translate
	s := Mat3Scale(2, 2)
	tr := Mat3Translate(3, 4)
	st := tr.Mul(s)
	v := st.MulVec2(NewVec2(1, 1))
	require.InDelta(t, 5.0, v.X, 1e-9)
	require.InDelta(t, 6.0, v.Y, 1e-9)
}

func TestMat3MulVec2(t *testing.T) {
	m := Mat3Translate(1, 2).Mul(Mat3Scale(3, 4))
	v := m.MulVec2(NewVec2(1, 1))
	require.InDelta(t, 4.0, v.X, 1e-9)
	require.InDelta(t, 6.0, v.Y, 1e-9)
}

func TestMat3MulVec2Dir(t *testing.T) {
	m := Mat3Translate(100, 200)
	// Direction should not be affected by translation
	v := m.MulVec2Dir(NewVec2(1, 0))
	require.InDelta(t, 1.0, v.X, 1e-9)
	require.InDelta(t, 0.0, v.Y, 1e-9)
}

func TestMat3MulVec3(t *testing.T) {
	m := Mat3Scale(2, 3)
	v := m.MulVec3(NewVec3(1, 1, 1))
	require.InDelta(t, 2.0, v.X, 1e-9)
	require.InDelta(t, 3.0, v.Y, 1e-9)
	require.InDelta(t, 1.0, v.Z, 1e-9)
}

func TestMat3Determinant(t *testing.T) {
	d := Mat3Identity().Determinant()
	require.InDelta(t, 1.0, d, 1e-9)

	d = Mat3Scale(2, 3).Determinant()
	require.InDelta(t, 6.0, d, 1e-9)
}

func TestMat3Inverse(t *testing.T) {
	m := Mat3Translate(5, 10).Mul(Mat3Scale(2, 3))
	inv, ok := m.Inverse()
	require.True(t, ok)
	product := m.Mul(inv)
	require.True(t, product.ApproxEqual(Mat3Identity(), 1e-9))
}

func TestMat3InverseSingular(t *testing.T) {
	m := Mat3{} // zero matrix
	_, ok := m.Inverse()
	require.False(t, ok)
}

func TestMat3Transpose(t *testing.T) {
	m := Mat3FromRows(
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	)
	mt := m.Transpose()
	require.True(t, m.ApproxEqual(mt.Transpose(), 1e-9))
	require.InDelta(t, m.At(0, 1), mt.At(1, 0), 1e-9)
	require.InDelta(t, m.At(1, 2), mt.At(2, 1), 1e-9)
}

func TestMat3ApproxEqual(t *testing.T) {
	a := Mat3Identity()
	b := Mat3Identity()
	b[0] = 1.0000000001
	require.True(t, a.ApproxEqual(b, 1e-9))
	require.False(t, a.ApproxEqual(Mat3Scale(2, 2), 1e-9))
}

func TestMat3At(t *testing.T) {
	m := Mat3Identity()
	require.InDelta(t, 1.0, m.At(0, 0), 1e-9)
	require.InDelta(t, 0.0, m.At(0, 1), 1e-9)
}

func TestMat3Set(t *testing.T) {
	m := Mat3Identity()
	m = m.Set(0, 1, 7.0)
	require.InDelta(t, 7.0, m.At(0, 1), 1e-9)
}

func TestMat3Col(t *testing.T) {
	m := Mat3Translate(5, 10)
	col2 := m.Col(2)
	require.InDelta(t, 5.0, col2.X, 1e-9)
	require.InDelta(t, 10.0, col2.Y, 1e-9)
	require.InDelta(t, 1.0, col2.Z, 1e-9)
}

func TestMat3Row(t *testing.T) {
	m := Mat3Translate(5, 10)
	row0 := m.Row(0)
	require.InDelta(t, 1.0, row0.X, 1e-9)
	require.InDelta(t, 0.0, row0.Y, 1e-9)
	require.InDelta(t, 5.0, row0.Z, 1e-9)
}

func TestMat3String(t *testing.T) {
	s := Mat3Identity().String()
	require.Contains(t, s, "Mat3[")
}

func TestMat3FromRows(t *testing.T) {
	m := Mat3FromRows(
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	)
	require.InDelta(t, 1.0, m.At(0, 0), 1e-9)
	require.InDelta(t, 2.0, m.At(0, 1), 1e-9)
	require.InDelta(t, 3.0, m.At(0, 2), 1e-9)
	require.InDelta(t, 4.0, m.At(1, 0), 1e-9)
	require.InDelta(t, 9.0, m.At(2, 2), 1e-9)
}
