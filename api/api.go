// Package api is a collection of open source golang projects repurposed or simplified into goFish functions and/or packages
package api

import (
	"github.com/edotau/goFish/simpleio"
	"github.com/panjf2000/ants"
	"github.com/pkg/browser"
)

func OpenDirectory(dir string) {
	err := browser.OpenFile(dir)
	simpleio.ErrorHandle(err)
}

func NewPool(size int) (*ants.Pool, error) {
	return ants.NewPool(size)

}
