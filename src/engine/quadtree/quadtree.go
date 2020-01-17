// Quadtree provides a data structure allowing efficient spatial queries on 2D points.
package quadtree

import (
	"INServer/src/common/logger"
	"INServer/src/engine/extensions/rect"
	"INServer/src/proto/engine"
	"errors"
)

// Quadtree is a dynamically resizable struct that offers
// performant access to items that can be assigned a 2d position. Internally,
// it represents the data with a tree by segmenting each level into four quadrants.
// The Quadtree has two attributes that may be tweaked depending on your use case:
//
// Maximum depth: To get the spatial resolution you want, the number of levels
// can be adjusted. For a depth D, the space will be subdivided 2^D times. The
// default is 10.
//
// Max items per node: The maximum number of items to store in a single node before
// it is split into smaller nodes. Note that if a node is full, and the maximum depth
// has been reached, max items per node will be ignored.
//
// The coordinate system of a Quadtree has origo at the lower left corner, with X
// and Z growing positively to the upper right corner. X is the horizontal axis,
// Z is the vertical axis.
//
// Quadtree is not thread-safe.
type Quadtree struct {
	maxDepth        int
	maxItemsPerNode int
	root            *node
	// The total number of items in this Quadtree.
	size            int
	debugAssertions bool
}

// An open-ended rectangle, where X, Z are inclusive
// but exclusive X+Width and Z+Height.

const quadrantNone = -1
const quadrantUpperLeft = 0
const quadrantUpperRight = 1
const quadrantLowerLeft = 2
const quadrantLowerRight = 3

// A single node in the tree.
type node struct {
	// The bounds that define this node
	bounds *engine.Rect
	// The depth this node is at. Root node is at depth 0.
	depth int
	// The entries that are inside this node
	items []treeEntry
	// The four child nodes of this node (one node per quadrant).
	ul, ur, ll, lr *node
}

// An individual entry in the tree.
type treeEntry struct {
	position *engine.Vector2
	uuid     string
}

// A function that takes treeEntries to iterate over entries.
type consumer func(treeEntry) bool

// Returns a Quadtree.
func NewQuadtree(bounds *engine.Rect, maxDepth, maxItemsPerNode int) (*Quadtree, error) {
	if maxDepth <= 0 {
		return nil, logger.Error("Creating tree failed: maxDepth must be larger than 1")
	}
	if maxItemsPerNode <= 0 {
		return nil, logger.Error("Creating tree failed: maxItemsPerNode must be larger than 0")
	}
	return &Quadtree{
		maxDepth:        maxDepth,
		maxItemsPerNode: maxItemsPerNode,
		root: &node{
			bounds: bounds,
			depth:  0,
			items:  make([]treeEntry, 0, 4),
		},
	}, nil
}

// Returns the number of items added to this tree.
func (qt *Quadtree) Size() int {
	return qt.size
}

// Returns the objects within the given bounds.
func (qt *Quadtree) Query(bounds *engine.Rect) []interface{} {
	items := make([]interface{}, 0, 10)
	queryInternal(qt.root, bounds, func(item treeEntry) bool {
		items = append(items, item.uuid)
		return true
	})
	return items
}

// Iterates over the items inside the given bounds. it will be invoked with each
// found data and its corresponding position until all items inside bounds have
// been found. The ordering is undefined.
func (qt *Quadtree) QueryIterative(bounds *engine.Rect, it func(string, *engine.Vector2) bool) {
	queryInternal(qt.root, bounds, func(item treeEntry) bool {
		return it(item.uuid, item.position)
	})
}

// This will recurse down the tree, removing the nodes that
// have no overlap with the given bounds. When all overlapping
// nodes are found, their items are returned.
// Returns false when search should be terminated.
func queryInternal(node *node, bounds *engine.Rect, consumer consumer) bool {
	if overlaps(node.bounds, bounds) {
		if node.items == nil {
			// This node has no items, but it has children. Keep recursing.
			keepGoing := queryInternal(node.ul, bounds, consumer)
			if keepGoing {
				keepGoing = queryInternal(node.ur, bounds, consumer)
			}
			if keepGoing {
				keepGoing = queryInternal(node.ll, bounds, consumer)
			}
			if keepGoing {
				queryInternal(node.lr, bounds, consumer)
			}
		} else {
			// We reached an end node. Since this node may only be partially
			// overlapping, ensure each item is inside bounds before consuming.
			for _, e := range node.items {
				if rect.Contains(bounds, e.position) {
					// If consumer returns false, stop immediately and terminate.
					if !consumer(e) {
						return false
					}
				}
			}
		}
	}
	return true
}

// Adds the data to the tree with the given position.
func (qt *Quadtree) Add(uuid string, position *engine.Vector2) (err error) {
	if !rect.Contains(qt.root.bounds, position) {
		// If the root node can't contain this data, signal error.
		return errors.New("Add failed: position outside bounds of tree.")
	}
	item := treeEntry{position, uuid}
	addInternal(qt, qt.root, item)
	qt.size += 1
	return err
}

// Recurses down the tree to find the correct node for the item.
func addInternal(qt *Quadtree, node *node, item treeEntry) {
	quadrant := whichQuadrant(node.bounds, item.position)
	if quadrant == quadrantNone {
		// The position doesn't belong to this node at all - exit.
		return
	}

	if node.items != nil && node.depth >= qt.maxDepth {
		// We've reached the max depth of the tree. The item must be stored
		// inside this node, regardless of maxItemsPerNode.
		node.items = append(node.items, item)
	} else if len(node.items) >= qt.maxItemsPerNode {
		// This node is already at max capacity, so we need to split it into
		// child nodes.
		ul, ur, ll, lr := rect.Quadrants(node.bounds)
		node.ul = newNode(node, ul)
		node.ur = newNode(node, ur)
		node.ll = newNode(node, ll)
		node.lr = newNode(node, lr)
		items := node.items
		node.items = nil
		for _, i := range items {
			addInternal(qt, node, i)
		}
		addInternal(qt, node, item)
	} else if node.items != nil {
		// This node still has an items array which means it does not have
		// any child nodes. Just append to this node.
		node.items = append(node.items, item)
	} else {
		// Recurse down the tree to find the correct node.
		switch quadrant {
		case quadrantUpperLeft:
			addInternal(qt, node.ul, item)
		case quadrantUpperRight:
			addInternal(qt, node.ur, item)
		case quadrantLowerLeft:
			addInternal(qt, node.ll, item)
		case quadrantLowerRight:
			addInternal(qt, node.lr, item)
		}
	}
}

func newNode(parent *node, bounds *engine.Rect) *node {
	return &node{
		bounds: bounds,
		depth:  parent.depth + 1,
		items:  make([]treeEntry, 0, 4),
	}
}

// Returns which quadrant the Point p is inside Rect r
func whichQuadrant(r *engine.Rect, p *engine.Vector2) int {
	if !rect.Contains(r, p) {
		return quadrantNone
	}
	midX := r.X + r.Width/2.0
	midZ := r.Z + r.Height/2.0
	if p.X < midX {
		// Left half
		if p.Z < midZ {
			return quadrantLowerLeft
		} else {
			return quadrantUpperLeft
		}
	} else {
		// Right half
		if p.Z < midZ {
			return quadrantLowerRight
		} else {
			return quadrantUpperRight
		}
	}
}

// Returns whether two rectangles overlap.
func overlaps(r1, r2 *engine.Rect) bool {
	return r1.X < r2.X+r2.Width &&
		r1.X+r1.Width > r2.X &&
		r1.Z+r1.Height > r2.Z &&
		r1.Z < r2.Z+r2.Height
}
