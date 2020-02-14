package vector3

import (
	"INServer/src/proto/engine"
	"math"
)

// Equal 判断相等
func Equal(a *engine.Vector3, b *engine.Vector3) bool {
	return a.X == b.X && a.Z == b.Z
}

// Normalize 归一化
func Normalize(a *engine.Vector3) *engine.Vector3 {
	length := math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
	return &engine.Vector3{
		X: a.X / length,
		Y: a.Y / length,
		Z: a.Z / length,
	}
}

// Multiply 乘法
func Multiply(a *engine.Vector3, scale float64) *engine.Vector3 {
	return &engine.Vector3{
		X: a.X * scale,
		Y: a.Y * scale,
		Z: a.Z * scale,
	}
}

// Minus 减法
func Minus(a *engine.Vector3, b *engine.Vector3) *engine.Vector3 {
	return &engine.Vector3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

// Add 加法
func Add(a *engine.Vector3, b *engine.Vector3) *engine.Vector3 {
	return &engine.Vector3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

// Dot 点积
func Dot(a *engine.Vector3, b *engine.Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}
