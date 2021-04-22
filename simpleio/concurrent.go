package simpleio

/*
import (
	"fmt"
	"os"
	"sync"
	"bytes"
	"log"
)

type chunk struct {
	bufsize int
	offset  int64
}

func ConcurrentReader(filename string) *bytes.Buffer {
	answer := bytes.Buffer{}
	const BufferSize = 100
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	filesize := int(fileinfo.Size())
	// Number of go routines we need to spawn.
	concurrency := filesize / BufferSize
	// buffer sizes that each of the go routine below should use. ReadAt
	// returns an error if the buffer size is larger than the bytes returned
	// from the file.
	chunksizes := make([]chunk, concurrency)

	// All buffer sizes are the same in the normal case. Offsets depend on the
	// index. Second go routine should start at 100, for example, given our
	// buffer size of 100.
	for i := 0; i < concurrency; i++ {
		chunksizes[i].bufsize = BufferSize
		chunksizes[i].offset = int64(BufferSize * i)
	}

	// check for any left over bytes. Add the residual number of bytes as the
	// the last chunk size.
	if remainder := filesize % BufferSize; remainder != 0 {
		c := chunk{bufsize: remainder, offset: int64(concurrency * BufferSize)}
		concurrency++
		chunksizes = append(chunksizes, c)
	}

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(chunksizes []chunk, i int) {
			defer wg.Done()
			chunk := chunksizes[i]
			buffer := make([]byte, chunk.bufsize)
			_, err := file.ReadAt(buffer, chunk.offset)
			if err != nil {
				log.Fatal(err)
			}
			//answer.Write(buffer)
			//fmt.Println("bytes read, string(bytestream): ", bytesread)
			fmt.Printf("%s", string(buffer))
		}(chunksizes, i)
	}

	wg.Wait()
	return &answer
}*/
