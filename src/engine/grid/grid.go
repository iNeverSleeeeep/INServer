package grid

import "INServer/src/proto/engine"

type (
	Item struct {
		uuid     string
		position *engine.Vector3
	}
	Grid struct {
		grids    [][]*Item
		gridSize float32
		width    int32
		height   int32
	}
)

func New(gridSize float32, width int32, height int32) *Grid {
	g := new(Grid)
	g.grids = make([][]*Item, width*height)
	g.gridSize = gridSize
	g.width = width
	g.height = height
	return g
}

func (g *Grid) Add(uuid string, position *engine.Vector3) {
	x := int32(position.X / g.gridSize)
	z := int32(position.Z / g.gridSize)
	index := g.getGirdIndex(x, z)
	if g.grids[index] == nil {
		g.grids[index] = make([]*Item, 0)
	}
	g.grids[index] = append(g.grids[index], &Item{uuid, position})
}

func (g *Grid) Remove(uuid string, position *engine.Vector3) {
	x := int32(position.X / g.gridSize)
	z := int32(position.Z / g.gridSize)
	index := g.getGirdIndex(x, z)
	if g.grids[index] == nil {
		g.grids[index] = make([]*Item, 0)
	}
	for i, item := range g.grids[index] {
		if item.uuid == uuid {
			g.grids[index] = append(g.grids[index][:i], g.grids[index][i+1:]...)
			break
		}
	}
}

func (g *Grid) Move(uuid string, from *engine.Vector3, to *engine.Vector3) {
	fromX := int32(from.X / g.gridSize)
	fromZ := int32(to.Z / g.gridSize)
	toX := int32(to.X / g.gridSize)
	toZ := int32(to.Z / g.gridSize)
	if fromX == toX && fromZ == toZ {
		return
	}
	g.Remove(uuid, from)
	g.Add(uuid, to)
}

func (g *Grid) GetNearItems(center *engine.Vector3) []*Item {
	x := int32(center.X / g.gridSize)
	z := int32(center.Z / g.gridSize)
	items := make([]*Item, 0)
	items = append(items, g.getGridItems(x, z)...)
	items = append(items, g.getGridItems(x-1, z-1)...)
	items = append(items, g.getGridItems(x-1, z)...)
	items = append(items, g.getGridItems(x-1, z+1)...)
	items = append(items, g.getGridItems(x, z-1)...)
	items = append(items, g.getGridItems(x, z+1)...)
	items = append(items, g.getGridItems(x+1, z-1)...)
	items = append(items, g.getGridItems(x+1, z)...)
	items = append(items, g.getGridItems(x+1, z+1)...)
	return items
}

func (g *Grid) getGirdIndex(x int32, z int32) int32 {
	return z*g.width + x
}

func (g *Grid) getGridItems(x int32, z int32) []*Item {
	if x >= 0 && z >= 0 {
		index := g.getGirdIndex(x, z)
		if index < int32(len(g.grids)) {
			return g.grids[index]
		}
	}
	return make([]*Item, 0)
}
