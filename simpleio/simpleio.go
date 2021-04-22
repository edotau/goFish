package simpleio

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	defaultBufSize = 4096
)

const mb = 1024 * 1024
const gb = 1024 * mb
const BufferSize = mb

// SimpleReader implements the io.Reader interface by providing
// the Read(b []byte) method. The struct contains an embedded *bufio.Reader
// and a pointer to os.File for closeure when reading is complete.
type SimpleReader struct {
	*bufio.Reader
	line   []byte
	Buffer *bytes.Buffer
	close  func() error
}

// Read reads data into p and is a method required to implement the io.Reader interface.
// It returns the number of bytes read into p.
func (reader *SimpleReader) Read(b []byte) (n int, err error) {
	return reader.Read(b)
}

// NewSimpleReader will process a given file and performs error handling if an error occurs.
// SimpleReader will prcoess gzipped files accordinging by performing a check on the suffix
// of the provided file.
func NewReader(filename string) *SimpleReader {
	if strings.HasPrefix(filename, "http") {
		return HttpReader(filename)
	}
	file := OpenFile(filename)
	var answer SimpleReader = SimpleReader{
		line:   make([]byte, defaultBufSize),
		Buffer: &bytes.Buffer{},
		close:  file.Close,
	}
	switch true {
	case strings.HasSuffix(filename, ".gz"):
		gzipReader, err := gzip.NewReader(file)
		ErrorHandle(err)
		answer.Reader = bufio.NewReader(gzipReader)
	default:
		answer.Reader = bufio.NewReader(file)
	}
	return &answer
}

// ReadLine will return a bytes.Buffer pointing to the internal slice of bytes. Provided this function is called within a loop,
// the function will read one line at a time, and return bool to continue reading. Important to note the buffer return points to
// the internal slice belonging to the reader, meaning the slice will be overridden if the data is not copied. Please be aware the
// reader will call close on the file once the reader encounters EOF.
func ReadLine(reader *SimpleReader) (*bytes.Buffer, bool) {
	var err error
	reader.line, err = reader.ReadSlice('\n')
	reader.Buffer.Reset()
	if err == nil {
		if reader.line[len(reader.line)-1] == '\n' {
			return BytesToBuffer(reader), false
		} else {
			log.Fatalf("Error: end of line did not end with an end of line character...\n")
		}
	} else {
		if err == bufio.ErrBufferFull {
			reader.line = readMore(reader)
			return BytesToBuffer(reader), false
		} else {
			CatchErrThrowEOF(err)
		}
	}
	return nil, true
}

// readMore is a private helper function to deal with very long lines to
// avoid alocating too much memory upfront and only resize the size of the buffer
// only when necessary.
func readMore(reader *SimpleReader) []byte {
	_, err := reader.Buffer.Write(reader.line)
	ErrorHandle(err)
	reader.line, err = reader.ReadSlice('\n')
	if err == nil {
		return reader.line
	}
	if err == bufio.ErrBufferFull {
		_, err = reader.Buffer.Write(reader.line)
		ErrorHandle(err)
		// recursive call to read next bytes until reaching end of line character
		return readMore(reader)
	}
	ErrorHandle(err)
	return reader.line
}

// BytesToBuffer will parse []byte and return a pointer to the same underlying bytes.Buffer
func BytesToBuffer(reader *SimpleReader) *bytes.Buffer {
	_, err := reader.Buffer.Write(reader.line[:len(reader.line)-1])
	ErrorHandle(err)
	return reader.Buffer
}

// CatchErrThrowEOF will silently handles and throws the EOF error and will log and exit any other errors.
func CatchErrThrowEOF(err error) {
	if err == io.EOF {
		return
	} else {
		ErrorHandle(err)
	}
}

// Close closes the File, rendering it unusable for I/O. On files that support SetDeadline,
// any pending I/O operations will be canceled and return immediately with an error.
// Close will return an error if it has already been called.
func (reader *SimpleReader) Close() {
	if reader != nil {
		ErrorHandle(reader.close())
	}
}

func GetBuffer(bufferPool sync.Pool) bytes.Buffer {
	return bufferPool.Get().(bytes.Buffer)
}

func PutBuffer(buf bytes.Buffer, bufferPool sync.Pool) {
	buf.Reset()
	bufferPool.Put(buf)
}

func simpleLogErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ErrorHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
func SimpleBufioPoolTest(filename string) {
	reader := NewSimpleReader(filename)
	var simplePool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	writer := fileio.EasyCreate("testdata/simplePoolTest.vcf.gz")

	var line []byte
	var err error
	var errW error
	var curr *bytes.Buffer
	for {
		curr = simplePool.Get().(*bytes.Buffer)

		line, err = reader.ReadSlice('\n')

		if err == nil {
			//line.b[len(line.b)-1] == '\n'
			curr.Write(line[:len(line)-1])

			//_, errW = fmt.Fprintf(writer, "%s\n", line)
			//fmt.Printf("%s\n", )
			simpleLogErr(errW)

			simplePool.Put(line)

		} else if err == io.EOF {
			CatchErrThrowEOF(err)
			break
		} else {
			simpleLogErr(err)
		}

	}
	writer.Close()
}*/

func OpenFile(filename string) *os.File {
	file, err := os.Open(filename)
	simpleLogErr(err)
	return file
}

// BufferPool implements a pool of bytes.Buffers in the form of a bounded
// channel.
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates a new BufferPool bounded to the given size.
func NewBufferPool(size int) *BufferPool {
	var bp BufferPool
	bp.pool.New = func() interface{} {
		return new(bytes.Buffer)
	}
	return &bp
}

// Get gets a Buffer from the BufferPool, or creates a new one if none are
// available in the pool.
func (bp *BufferPool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

// Put returns the given Buffer to the BufferPool.
func (bp *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	bp.pool.Put(b)
}

/*
func client() {
    for {
        var b *bytes.Buffer
        // Grab a buffer if available; allocate if not.
        select {
        case b = <-freeList:
            // Got one; nothing more to do.
        default:
            // None free, so allocate a new one.
            b = new(Buffer)
        }
        load(b)              // Read next message from the net.
        serverChan <- b      // Send to server.
    }
}

func server() {
    for {
        b := <-serverChan    // Wait for work.
        process(b)
        // Reuse buffer if there's room.
        select {
        case freeList <- b:
            // Buffer on free list; nothing more to do.
        default:
            // Free list full, just carry on.
        }
    }
}*/
/*
func routineReader(reader *SimpleReader, data chan []byte) {
	for curr, done := ReadLine(reader); !done; curr, done = ReadLine(reader) {
		data <- curr
	}
	close(data)
}

func ConcurrentReader(filename string) <- chan []byte {
	simpleReader := NewSimpleReader(filename)
	ans := make(chan []byte, 100000)
	go routineReader(simpleReader, ans)
	return ans
}*/
