package gamemap

import (
	"INServer/src/common/logger"
	"INServer/src/engine/quadtree"
	"INServer/src/proto/config"
	"INServer/src/proto/data"
)

type (
	Scene struct {
		masterMap *Map
		tree      *quadtree.Quadtree
	}
)

func NewScene(masterMap *Map, sceneConfig *config.Scene) *Scene {
	tree, err := quadtree.NewQuadtree(sceneConfig.Rect, 10, 1000)
	if err != nil {
		logger.Error(err)
		return nil
	}
	s := new(Scene)
	s.masterMap = masterMap
	s.tree = tree
	return s
}

func (s *Scene) EntityEnter(uuid string, entity *gameentity.Entity) {

}

func (s *Scene) EntityLeave(uuid string, entity *gameentity.Entity) {

}
