package vector3

import (
	"INServer/src/proto/engine"
	"math"
)

func Equal(a *engine.Vector3, b *engine.Vector3) bool {
	return a.X == b.X && a.Z == b.Z
}

func Normalize(a *engine.Vector3) *engine.Vector3 {
	length := math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
	return &engine.Vector3{
		X: a.X / length,
		Y: a.Y / length,
		Z: a.Z / length,
	}
}

func Multiply(a *engine.Vector3, scale float64) *engine.Vector3 {
	return &engine.Vector3{
		X: a.X * scale,
		Y: a.Y * scale,
		Z: a.Z * scale,
	}
}

func Minus(a *engine.Vector3, b *engine.Vector3) *engine.Vector3 {
	return &engine.Vector3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}
