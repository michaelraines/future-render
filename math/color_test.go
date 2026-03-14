package math

import "testing"

func TestColorFromRGBA(t *testing.T) {
	c := ColorFromRGBA(255, 128, 0, 255)
	if !ApproxEqual(c.R, 1, 0.01) || !ApproxEqual(c.A, 1, testEpsilon) {
		t.Errorf("unexpected color: %v", c)
	}
	if !ApproxEqual(c.G, 128.0/255.0, 0.01) {
		t.Errorf("unexpected green: %g", c.G)
	}
}

func TestColorFromHex(t *testing.T) {
	c := ColorFromHex(0xFF0000FF)
	if !c.ApproxEqual(ColorRed, testEpsilon) {
		t.Errorf("expected red, got %v", c)
	}
}

func TestColorRGBA(t *testing.T) {
	r, g, b, a := ColorWhite.RGBA()
	if r != 255 || g != 255 || b != 255 || a != 255 {
		t.Errorf("expected (255,255,255,255), got (%d,%d,%d,%d)", r, g, b, a)
	}
}

func TestColorLerp(t *testing.T) {
	c := ColorBlack.Lerp(ColorWhite, 0.5)
	if !ApproxEqual(c.R, 0.5, testEpsilon) || !ApproxEqual(c.G, 0.5, testEpsilon) {
		t.Errorf("expected 0.5 components, got %v", c)
	}
}

func TestColorPremultiply(t *testing.T) {
	c := NewColor(1, 0.5, 0.25, 0.5).Premultiply()
	if !ApproxEqual(c.R, 0.5, testEpsilon) || !ApproxEqual(c.G, 0.25, testEpsilon) {
		t.Errorf("unexpected premultiplied color: %v", c)
	}
	if !ApproxEqual(c.A, 0.5, testEpsilon) {
		t.Errorf("alpha should be unchanged: %g", c.A)
	}
}

func TestColorSRGBRoundtrip(t *testing.T) {
	c := NewColor(0.5, 0.2, 0.8, 1.0)
	roundtrip := c.ToSRGB().ToLinear()
	if !c.ApproxEqual(roundtrip, 1e-6) {
		t.Errorf("sRGB roundtrip failed: %v -> %v", c, roundtrip)
	}
}

func TestColorMatrixIdentity(t *testing.T) {
	cm := ColorMatrixIdentity()
	c := NewColor(0.3, 0.6, 0.9, 1.0)
	result := cm.Apply(c)
	if !result.ApproxEqual(c, testEpsilon) {
		t.Errorf("identity should not change color: %v -> %v", c, result)
	}
}

func TestColorMatrixScale(t *testing.T) {
	cm := ColorMatrixScale(0.5, 2, 0, 1)
	c := NewColor(1, 0.5, 1, 1)
	result := cm.Apply(c)
	expected := NewColor(0.5, 1, 0, 1)
	if !result.ApproxEqual(expected, testEpsilon) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestColorMatrixConcat(t *testing.T) {
	a := ColorMatrixScale(2, 1, 1, 1)
	b := ColorMatrixTranslate(0.1, 0, 0, 0)
	combined := a.Concat(b) // a applied after b
	c := NewColor(0.5, 0.5, 0.5, 1)
	step := b.Apply(c)
	step = a.Apply(step)
	direct := combined.Apply(c)
	if !step.ApproxEqual(direct, testEpsilon) {
		t.Errorf("concat mismatch: step=%v, direct=%v", step, direct)
	}
}
