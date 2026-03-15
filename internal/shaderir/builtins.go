package shaderir

// builtinFuncs maps Kage built-in function names to their GLSL equivalents.
// Functions with the same name in both languages are included for completeness.
var builtinFuncs = map[string]string{
	// Math functions (same name in Kage and GLSL).
	"sin":         "sin",
	"cos":         "cos",
	"tan":         "tan",
	"asin":        "asin",
	"acos":        "acos",
	"atan":        "atan",
	"atan2":       "atan", // Kage atan2(y,x) → GLSL atan(y,x)
	"pow":         "pow",
	"exp":         "exp",
	"log":         "log",
	"exp2":        "exp2",
	"log2":        "log2",
	"sqrt":        "sqrt",
	"inversesqrt": "inversesqrt",
	"abs":         "abs",
	"sign":        "sign",
	"floor":       "floor",
	"ceil":        "ceil",
	"fract":       "fract",
	"mod":         "mod",
	"min":         "min",
	"max":         "max",
	"clamp":       "clamp",
	"mix":         "mix",
	"step":        "step",
	"smoothstep":  "smoothstep",
	"length":      "length",
	"distance":    "distance",
	"dot":         "dot",
	"cross":       "cross",
	"normalize":   "normalize",
	"faceforward": "faceforward",
	"reflect":     "reflect",
	"refract":     "refract",
	"transpose":   "transpose",
	"dfdx":        "dFdx",
	"dfdy":        "dFdy",
	"fwidth":      "fwidth",

	// Ebitengine-specific functions that map to texture lookups.
	// These are handled specially by the transpiler.
}

// imageBuiltins lists image-related built-in functions that need special
// handling during transpilation.
var imageBuiltins = map[string]bool{
	"imageSrc0At":         true,
	"imageSrc1At":         true,
	"imageSrc2At":         true,
	"imageSrc3At":         true,
	"imageSrc0UnsafeAt":   true,
	"imageSrc1UnsafeAt":   true,
	"imageSrc2UnsafeAt":   true,
	"imageSrc3UnsafeAt":   true,
	"imageSrc0Origin":     true,
	"imageSrc1Origin":     true,
	"imageSrc2Origin":     true,
	"imageSrc3Origin":     true,
	"imageSrc0Size":       true,
	"imageSrc1Size":       true,
	"imageSrc2Size":       true,
	"imageSrc3Size":       true,
	"imageDstOrigin":      true,
	"imageDstSize":        true,
	"imageDstTextureSize": true,
}

// IsBuiltinFunc returns true if name is a Kage built-in function.
func IsBuiltinFunc(name string) bool {
	if _, ok := builtinFuncs[name]; ok {
		return true
	}
	return imageBuiltins[name]
}

// GLSLBuiltin returns the GLSL equivalent of a Kage built-in function name.
// For image built-ins, returns empty string (they need special handling).
func GLSLBuiltin(name string) string {
	return builtinFuncs[name]
}

// IsImageBuiltin returns true if name is a Kage image-related built-in.
func IsImageBuiltin(name string) bool {
	return imageBuiltins[name]
}

// IsConstructor returns true if name is a Kage type constructor (vec2, mat4, etc.).
func IsConstructor(name string) bool {
	_, err := ParseType(name)
	return err == nil
}
