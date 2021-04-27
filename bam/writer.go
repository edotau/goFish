package bam

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"log"
	"os"

	"github.com/edotau/goFish/simpleio"
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

/*
type BAMWriters struct {
	BarcodeSortedBam     *BAMWriter
	PositionBucketedBams map[string][]*BAMWriter
	positionChunkSize    int
	debugTags            bool
	channel              chan *Data
	// This mutex is Rlocked by the worker thread. When we close, we wait
	 //for the mutex to be unlocked to ensure data is flushed before continueing

	done sync.RWMutex
}

type BAMWriter struct {
	Writer  *bam.Writer
	Contigs map[string]*bam.Reference
	Record  bam.Record
}

func CreateBAM(ref *gobwa.GoBwaReference, path, read_groups, sample_id string, firstChunk bool) (*BAMWriter, error) {
	bw := &BAMWriter{}
	bw.Contigs = make(map[string]*bam.Reference)

	references := make([]*bam.Reference, 0, 0)

	gobwa.EnumerateContigs(ref, func(name string, length int) {
		r, err := bam.NewReference(name, name, "human", length, nil, nil)
		if err != nil {
			panic(err)
		}
		references = append(references, r)
		bw.Contigs[name] = r
	})

	// Only include the CO headers on the first chunk: avoid having them duplicated during samtools merge
	comments := []byte("")
	if firstChunk {
		comments = []byte("@CO\t10x_bam_to_fastq:R1(RX:QX,TR:TQ,SEQ:QUAL)\n@CO\t10x_bam_to_fastq:R2(SEQ:QUAL)\n@CO\t10x_bam_to_fastq:I1(BC:QT)")
	}
	h, err := bam.NewHeader(comments, references)

	if err != nil {
		panic(err)
	}

	// NewReadGroup(name, center, desc, lib, prog, plat, unit, sample string, date time.Time, size int, flow, key []byte)
	for _, rg_id := range strings.Split(read_groups, ",") {
		// currently, the ID is composed of:
		// sample:library:gem_group:flowcell:lane
		rg_fields := strings.Split(rg_id, ":")
		if len(rg_fields) == 0 {
			log.Printf("Empty RG was specified, skipping")
		} else if len(rg_fields) < 5 {
			log.Printf("RG is not fully specified, skipping: %s", rg_id)
		} else {
			rg, err := bam.NewReadGroup(
				rg_id, //ID
				"",    //CN
				"",    //DS
				rg_fields[1]+"."+rg_fields[2], //LB = (input library).(gem group)
				"",           //PG
				"ILLUMINA",   //PL
				rg_id,        //PU: just make same as ID?
				rg_fields[0], //SM
				time.Now(),
				0,
				nil,
				nil)
			if err != nil {
				panic(err)
			}
			h.AddReadGroup(rg)
		}
	}

	// Add a program line for lariat
	prog := bam.NewProgram(
		"lariat",                   // ID
		"longranger.lariat",        // PN
		strings.Join(os.Args, " "), // CL
		"",          // PP - no need to indicate previous, since Lariat produces the initial BAM
		__VERSION__) // VN
	h.AddProgram(prog)

	file, err := os.Create(path)

	if err != nil {
		return nil, err
	}

	w, err := bam.NewWriter(file, h, 2)

	if err != nil {
		panic(err)
	}
	bw.Writer = w
	return bw, nil
}

*/
