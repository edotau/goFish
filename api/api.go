// Package api is a collection of open source golang projects repurposed or simplified into goFish functions and/or packages
package api

import (
	"github.com/edotau/goFish/simpleio"
	"github.com/pkg/browser"
)

func OpenDirectory(dir string) {
	err := browser.OpenFile(dir)
	simpleio.ErrorHandle(err)
}
