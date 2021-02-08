// Copyright 2021 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package prtree

import (
	"github.com/tidwall/geoindex/child"
	"github.com/tidwall/ptree"
	"github.com/tidwall/rtree"
)

// PRTree is a tree for storing points and rects
type PRTree struct {
	ptree *ptree.PTree
	rtree *rtree.RTree
}

// New returns a new PRTree
func New(min, max [2]float64) *PRTree {
	tr := new(PRTree)
	tr.ptree = ptree.New(min, max)
	tr.rtree = new(rtree.RTree)
	return tr
}

// Insert an item into the structure
func (tr *PRTree) Insert(min, max [2]float64, data interface{}) {
	if min == max && tr.ptree.InBounds(min) {
		tr.ptree.Insert(min, data)
	} else {
		tr.rtree.Insert(min, max, data)
	}
}

// Delete an item from the structure
func (tr *PRTree) Delete(min, max [2]float64, data interface{}) {
	if min == max && tr.ptree.InBounds(min) {
		tr.ptree.Delete(min, data)
	} else {
		tr.rtree.Delete(min, max, data)
	}
}

// Replace an item in the structure. This is effectively just a Delete
// followed by an Insert. But for some structures it may be possible to
// optimize the operation to avoid multiple passes
func (tr *PRTree) Replace(
	oldMin, oldMax [2]float64, oldData interface{},
	newMin, newMax [2]float64, newData interface{},
) {
	tr.Delete(oldMin, oldMax, oldData)
	tr.Insert(newMin, newMax, newData)
}

// Search the structure for items that intersects the rect param
func (tr *PRTree) Search(
	min, max [2]float64,
	iter func(min, max [2]float64, data interface{}) bool,
) {
	var quit bool
	tr.ptree.Search(min, max, func(point [2]float64, data interface{}) bool {
		if !iter(point, point, data) {
			quit = true
			return false
		}
		return true
	})
	if !quit {
		tr.rtree.Search(min, max, iter)
	}
}

// Scan iterates through all data in tree in no specified order.
func (tr *PRTree) Scan(iter func(min, max [2]float64, data interface{}) bool) {
	var quit bool
	tr.ptree.Scan(func(point [2]float64, data interface{}) bool {
		if !iter(point, point, data) {
			quit = true
			return false
		}
		return true
	})
	if !quit {
		tr.rtree.Scan(iter)
	}
}

// Len returns the number of items in tree
func (tr *PRTree) Len() int {
	return tr.ptree.Len() + tr.rtree.Len()
}

func expand(amin, amax, bmin, bmax [2]float64) (min, max [2]float64) {
	if bmin[0] < amin[0] {
		amin[0] = bmin[0]
	}
	if bmax[0] > amax[0] {
		amax[0] = bmax[0]
	}
	if bmin[1] < amin[1] {
		amin[1] = bmin[1]
	}
	if bmax[1] > amax[1] {
		amax[1] = bmax[1]
	}
	return amin, amax
}

// Bounds returns the minimum bounding box
func (tr *PRTree) Bounds() (min, max [2]float64) {
	if tr.ptree.Len() > 0 {
		amin, amax := tr.ptree.MinBounds()
		if tr.rtree.Len() > 0 {
			bmin, bmax := tr.rtree.Bounds()
			min, max = expand(amin, amax, bmin, bmax)
		} else {
			min, max = amin, amax
		}
	} else if tr.rtree.Len() > 0 {
		min, max = tr.rtree.Bounds()
	}
	return min, max
}

type pChildNode struct {
	data interface{}
}
type rChildNode struct {
	data interface{}
}

// Children returns all children for parent node. If parent node is nil
// then the root nodes should be returned.
// The reuse buffer is an empty length slice that can optionally be used
// to avoid extra allocations.
func (tr *PRTree) Children(parent interface{}, reuse []child.Child,
) (children []child.Child) {
	children = reuse[:0]

	switch parent := parent.(type) {
	case nil:
		// Gather the parent nodes for all PTree and RTree children.
		// For all children that are node data (not item data), we'll need to
		// wrap the data with a local type to keep track of the which tree the
		// node belongs to.
		children = tr.ptree.Children(nil, children)
		for i := 0; i < len(children); i++ {
			if !children[i].Item {
				children[i].Data = pChildNode{children[i].Data}
			}
		}
		mark := len(children)
		children = tr.rtree.Children(nil, children)
		for i := mark; i < len(children); i++ {
			if !children[i].Item {
				children[i].Data = rChildNode{children[i].Data}
			}
		}
		return children
	case pChildNode:
		children = tr.ptree.Children(parent.data, children)
		for i := 0; i < len(children); i++ {
			if !children[i].Item {
				children[i].Data = pChildNode{children[i].Data}
			}
		}
	case rChildNode:
		children = tr.rtree.Children(parent.data, children)
		for i := 0; i < len(children); i++ {
			if !children[i].Item {
				children[i].Data = rChildNode{children[i].Data}
			}
		}
	default:
		panic("invalid node type")
	}
	return children
}
