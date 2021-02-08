// Copyright 2021 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package prtree

import (
	"math/rand"
	"testing"
	"time"

	"github.com/tidwall/geoindex"
)

func init() {
	seed := time.Now().UnixNano()
	seed = 1612713107378415000
	println("seed:", seed)
	rand.Seed(seed)
}

func newWorld() *PRTree {
	return New([2]float64{-180, -90}, [2]float64{180, 90})
}

// IsMixedTree is needed for the geoindex.Tests.TestBenchVarious
// Do not remove this function.
func (tr *PRTree) IsMixedTree() bool {
	return true
}

func TestGeoIndex(t *testing.T) {
	t.Run("BenchVarious", func(t *testing.T) {
		geoindex.Tests.TestBenchVarious(t, newWorld(), 1000000)
	})
	t.Run("RandomRects", func(t *testing.T) {
		geoindex.Tests.TestRandomRects(t, newWorld(), 10000)
	})
	t.Run("RandomPoints", func(t *testing.T) {
		geoindex.Tests.TestRandomPoints(t, newWorld(), 10000)
	})
	t.Run("ZeroPoints", func(t *testing.T) {
		geoindex.Tests.TestZeroPoints(t, newWorld())
	})
	t.Run("CitiesSVG", func(t *testing.T) {
		geoindex.Tests.TestCitiesSVG(t, newWorld())
	})
}

func BenchmarkRandomInsert(b *testing.B) {
	geoindex.Tests.BenchmarkRandomInsert(b, newWorld())
}
