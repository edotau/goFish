package bam

import (
	"fmt"
	"github.com/edotau/goFish/simpleio"
	"net/http"
	"strings"
	"sync"
)

func ViewUrl(url string) {
	resp, err := http.Get(url)
	simpleio.FatalErr(err)
	ans := make(chan Sam)
	if strings.HasSuffix(url, ".bam") {
		reader := &BamReader{}
		//reader.File = resp.Body
		reader.Gunzip = NewBgzipReader(resp.Body)
		h := ReadHeader(reader)
		binaryData := make(chan *BinaryDecoder)

		var wg sync.WaitGroup
		go BamToChannel(reader, binaryData)

		go func() {
			for each := range binaryData {
				ans <- *BamBlockToSam(h, each)
			}
			wg.Done()
		}()
		go func() {
			wg.Wait()
			close(ans)
		}()
	}
	if strings.HasSuffix(url, ".sam") {
		reader := NewSamReader(url)
		go func() {
			for curr, done := UnmarshalSam(reader); !done; curr, done = UnmarshalSam(reader) {
				ans <- *curr
			}
			close(ans)
		}()
	}
	for i := range ans {
		fmt.Printf("%s\n", ToString(&i))
	}
}
