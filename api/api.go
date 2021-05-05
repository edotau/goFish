// Package api is a collection of open source golang projects repurposed or simplified into goFish functions and/or packages
package api

import (
	"github.com/edotau/goFish/simpleio"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/browser"
)

func OpenDirectory(dir string) {
	err := browser.OpenFile(dir)
	simpleio.FatalErr(err)
}

func NewPool(size int) *ants.Pool {
	pool, err := ants.NewPool(size)
	simpleio.FatalErr(err)
	return pool
}
