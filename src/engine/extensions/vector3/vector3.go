package vector3

import "INServer/src/proto/engine"

func Equal(a *engine.Vector3, b *engine.Vector3) bool {
	return a.X == b.X && a.Z == b.Z
}
