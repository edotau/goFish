package bam

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"github.com/goFish/simpleio"
	"log"
	"os"
)

var bamMagic = [4]byte{'B', 'A', 'M', 0x1}

type BamWriter struct {
	gzip.Writer
	File   *os.File
	Stream *bytes.Buffer
	buf    [4]byte
}

func NewBamWriter(filename string) *BamWriter {
	ans := &BamWriter{}
	ans.File = simpleio.OpenFile(filename)
	r := bufio.NewWriter(ans.File)
	//common.ExitIfError(err)
	ans.Writer = *gzip.NewWriter(r)
	return ans
}

func WriteBinaryHeader(filename string, bh *Header) *BamWriter {
	writer := NewBamWriter(filename)
	binary.Write(writer, binary.LittleEndian, bamMagic)
	text := MarshalText(bh)
	binary.Write(writer, binary.LittleEndian, int32(len(text)))
	writer.Write(text)
	binary.Write(writer, binary.LittleEndian, int32(len(bh.Chroms)))
	var name []byte
	for _, i := range bh.Chroms {
		name = append(name, []byte(i.Name)...)
		name = append(name, 0)
		binary.Write(writer, binary.LittleEndian, int32(len(name)))
		writer.Write(name)
		name = name[:0]
		binary.Write(writer, binary.LittleEndian, int32(i.Size))
	}
	writer.Flush()
	return writer
}

var (
	lenFieldSize      = binary.Size(BinaryDecoder{}.BlockSize)
	bamFixedRemainder = binary.Size(BinaryDecoder{}) - lenFieldSize
)

func WriteBam(writer *BamWriter, record *Sam) {
	if len(record.RName) > 254 {
		log.Fatalf("Error: length of reference name is too long...\n")
	}
	recLen := bamFixedRemainder + len(record.RName) + 1 + len(record.Cigar)<<2 + len(record.Seq)
	WriteInt32(int32(recLen), writer)

}

/*

func TestWriter(t *testing.T) {
	for _, test := range readBamTests {
		bamFile := NewBamReader(test.bam)
		header := ReadHeader(bamFile)
		WriteBinaryHeader("testdata/writer.bam", header)
		wBamFile := NewBamReader("testdata/writer.bam")
		readWriteHeader := ReadHeader(wBamFile)
		fmt.Printf("%s", readWriteHeader.Txt)
	}


}
*/

func WriteUint8(v uint8, writer *BamWriter) {
	writer.buf[0] = v
	writer.Write(writer.buf[:1])
}

func WriteUint16(v uint16, writer *BamWriter) {
	binary.LittleEndian.PutUint16(writer.buf[:2], v)
	writer.Write(writer.buf[:2])
}

func WriteInt32(v int32, writer *BamWriter) {
	binary.LittleEndian.PutUint32(writer.buf[:4], uint32(v))
	writer.Write(writer.buf[:4])
}

func WriteUint32(v uint32, writer *BamWriter) {
	binary.LittleEndian.PutUint32(writer.buf[:4], v)
	writer.Write(writer.buf[:4])
}

func MarshalText(bh *Header) []byte {
	return bh.Text.Bytes()
}
