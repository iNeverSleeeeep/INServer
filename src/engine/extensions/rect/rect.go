package rect

import "INServer/src/proto/engine"

func Quadrants(r *engine.Rect) (ul, ur, ll, lr *engine.Rect) {
	w := r.Width / 2.0
	h := r.Height / 2.0
	ll = &engine.Rect{X: r.X, Z: r.Z, Width: w, Height: h}
	ul = &engine.Rect{X: r.X, Z: r.Z + h, Width: w, Height: h}
	ur = &engine.Rect{X: r.X + w, Z: r.Z + h, Width: w, Height: h}
	lr = &engine.Rect{X: r.X + w, Z: r.Z, Width: w, Height: h}
	return ul, ur, ll, lr
}

func Contains(r *engine.Rect, p *engine.Vector2) bool {
	return p.X >= r.X && p.Z >= r.Z && p.X < r.X+r.Width && p.Z < r.Z+r.Height
}
