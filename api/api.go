package api

import (
	"github.com/goFish/simpleio"
	"github.com/pkg/browser"
)

func OpenDirectory(dir string) {
	err := browser.OpenFile(dir)
	simpleio.ErrorHandle(err)
}
