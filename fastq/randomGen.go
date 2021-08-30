package fastq

import (
	"fmt"

	"github.com/edotau/goFish/code"
	"github.com/edotau/goFish/fasta"
	"github.com/edotau/goFish/stats"
)

func RandomPairedReads(genome []fasta.Fasta, readLength int, numReads int, insertSize int) []PairedEnd {
	var randFqs []PairedEnd = make([]PairedEnd, numReads)
	for i := 0; i < numReads; i++ {
		randFqs[i] = RandChrSeq(genome, readLength, insertSize)
		fmt.Printf("%s\n", randFqs[i].ReadOne.ToString())
	}

	return randFqs
}

// TODO implement logic to randomize strand
func RandChrSeq(genome []fasta.Fasta, readLength int, insertSize int) PairedEnd {
	fq := PairedEnd{
		ReadOne: Fastq{},
		ReadTwo: Fastq{},
	}
	var chr int = stats.RandIntInRange(0, len(genome))
	var randPos int = stats.RandIntInRange(0, len(genome[chr].Seq))
	var start, end int

	if randPos+insertSize > len(genome[chr].Seq) {
		start = len(genome[chr].Seq) - insertSize + 1
		end = len(genome[chr].Seq)

		fq.ReadOne.Name = fmt.Sprintf("%s_%d_%d_%c_R:1", genome[chr].Name, start, end, '+')
		fq.ReadOne.Seq = append(fq.ReadOne.Seq, genome[chr].Seq[start:end]...)

		start = len(genome[chr].Seq) - readLength
		end = len(genome[chr].Seq)

		fq.ReadTwo.Name = fmt.Sprintf("%s_%d_%d_%c_R:2", genome[chr].Name, start, end, '-')
		fq.ReadTwo.Seq = append(fq.ReadTwo.Seq, genome[chr].Seq[start:end]...)

		fmt.Printf("%s\n", fq.ReadOne.ToString())
		return fq

	} else {
		start, end = randPos, randPos+insertSize
		fq.ReadOne.Name = fmt.Sprintf("%s_%d_%d_%c_R:1", genome[chr].Name, start, end, '+')
		fq.ReadOne.Seq = append(fq.ReadOne.Seq, genome[chr].Seq[start:end]...)

		start = len(genome[chr].Seq) - readLength
		end = len(genome[chr].Seq)

		fq.ReadTwo.Name = fmt.Sprintf("%s_%d_%d_%c_R:2", genome[chr].Name, start, end, '-')
		fq.ReadTwo.Seq = append(fq.ReadTwo.Seq, genome[chr].Seq[start:end]...)
		return fq
	}
}

func NewFastq(length int) Fastq {
	return Fastq{
		Name: "",
		Seq:  make([]code.Dna, length),
		Qual: GenerateFakeQual(length),
	}
}

func GenerateFakeQual(length int) []byte {
	var answer []byte = make([]byte, length)
	var asci = []byte{'!', '#', '$', '%', '&', '(', ')', '*', '+', '`', '-', '.', '/', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?', '@', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'}

	for i := 0; i < length; i++ {
		answer[i] = asci[stats.RandIntInRange(0, len(asci))]
	}
	return answer
}

func GenerateFakeBases(length int) []code.Dna {
	answer := make([]code.Dna, length)

	for i := 0; i < length; i++ {
		answer[i] = code.NoMaskDnaArray[stats.RandIntInRange(0, len(code.NoMaskDnaArray))]
	}
	return answer
}

func CreateEmptyFq() Fastq {
	fq := Fastq{Name: "null", Seq: nil}
	return fq
}