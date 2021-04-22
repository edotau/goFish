package main

import (
	"fmt"
	"strings"

	"github.com/goFish/bam"
	"github.com/goFish/simpleio"
	"github.com/goFish/vcf"
	"github.com/vertgenlab/gonomics/dna"
)

func getVcfSeq(v *vcf.Vcf) [][]dna.Base {
	ans := make([][]dna.Base, 0)
	ans = append([][]dna.Base{dna.StringToBases(v.Ref)}, GetAltBases(v.Alt)...)
	return ans
}

func GetAltBases(alt string) [][]dna.Base {
	words := strings.Split(alt, ",")
	var answer [][]dna.Base = make([][]dna.Base, len(words))
	for i := 0; i < len(words); i++ {
		answer[i] = dna.StringToBases(words[i])
	}
	return answer
}

type index struct {
	p1 int
	p2 int
	f1 int
}

type family struct {
	One vcf.Genotype
	Two vcf.Genotype
	F1  vcf.Genotype
	Seq [][]dna.Base
}

func buildGenotypeMap(v *vcf.Vcf, vcfHeader *vcf.Header, famIdx index, mapToVcf map[uint64]family) map[uint64]family {
	code := vcf.GenomeKeyUint64(vcfHeader.Ref[v.Chr], v.Pos-1)
	_, ok := mapToVcf[code]
	if !ok {
		mapToVcf[code] = getGenotypeData(v, famIdx)
	}
	return mapToVcf
}

func getGenotypeData(v *vcf.Vcf, famIdx index) family {
	ans := family{}
	ans.One = vcf.ParseGt(v.Genotypes[famIdx.p1])
	ans.Two = vcf.ParseGt(v.Genotypes[famIdx.p2])
	ans.F1 = vcf.ParseGt(v.Genotypes[famIdx.f1])
	ans.Seq = getVcfSeq(v)
	return ans
}

func SnpSearch(samfile string, genotypeVcf string, fOne string, parentOne string, parentTwo, prefix string) {
	vcfs := make(chan vcf.Vcf)

	//sampleHash := vcf.HeaderToMaps(reader.Header)
	snpDb := make(map[uint64]family)

	file := vcf.NewReader(genotypeVcf)
	vcfHeader := vcf.ReadHeader(file)

	go vcf.ReadToChan(file, vcfs)

	famIdx := index{p1: vcfHeader.Samples[parentOne], p2: vcfHeader.Samples[parentTwo], f1: vcfHeader.Samples[fOne]}
	for genotype := range vcfs {
		if vcf.ASFilter(&genotype, famIdx.p1, famIdx.p2, famIdx.f1) {
			buildGenotypeMap(&genotype, vcfHeader, famIdx, snpDb)
		}
	}

	childOne, childTwo := simpleio.NewWriter(fmt.Sprintf("%s.%s.SNPs.sam", prefix, parentOne)), simpleio.NewWriter(fmt.Sprintf("%s.%s.SNPs.sam", prefix, parentTwo))
	defer childOne.Close()
	defer childTwo.Close()

	header, sams := bam.Read(samfile)
	childOne.WriteString(header.Text.String())
	childOne.WriteByte('\n')
	childTwo.WriteString(header.Text.String())
	childTwo.WriteByte('\n')

	var i, parentAllele1, parentAllele2 int
	var target, query, j int
	var ok bool
	var code uint64

	var gV family
	//for read, done := bam.UnmarshalSam(samFile); done != true; read, done = bam.UnmarshalSam(samFile) {
	for read := range sams {
		parentAllele1, parentAllele2 = 0, 0
		target = read.Pos - 1
		query = 0
		for i = 0; i < len(read.Cigar); i++ {
			switch read.Cigar[i].Op {
			case 'S':
				query += int(read.Cigar[i].RunLen)
			case 'I':
				//TODO: Figure out how to take insertions into account. This algorithm below should work in theory, but there is a case I can't figure out
				//code = ChromPosToUInt64(int(vcfHeader.Ref[read.RName]), int(target))
				//_, ok = snpDb[code]
				//if ok {
				//	gV = snpDb[code]
				//	if dna.CompareSeqsIgnoreCase(read.Seq[query:query+read.Cigar[i].RunLen], gV.Alleles[vcf.GetGenotypes(gV.Genotypes)[vcfHeader.Samples[parentOne]].AlleleOne]) == 0 && dna.CompareSeqsIgnoreCase(read.Seq[query:query+read.Cigar[i].RunLen], gV.Alleles[vcf.GetGenotypes(gV.Genotypes)[vcfHeader.Samples[parentOne]].AlleleTwo]) == 0 {
				//		parentAllele1++
				//	}
				//	if dna.CompareSeqsIgnoreCase(read.Seq[query:query+read.Cigar[i].RunLen], gV.Alleles[vcf.GetGenotypes(gV.Genotypes)[vcfHeader.Samples[parentTwo]].AlleleOne]) == 0 && dna.CompareSeqsIgnoreCase(read.Seq[query:query+read.Cigar[i].RunLen], gV.Alleles[vcf.GetGenotypes(gV.Genotypes)[vcfHeader.Samples[parentTwo]].AlleleTwo]) == 0 {
				//		parentAllele2++
				//	}
				//}
				query += int(read.Cigar[i].RunLen)
			case 'D':
				code = vcf.GenomeKeyUint64(int(vcfHeader.Ref[read.RName]), int(target))
				gV, ok = snpDb[code]
				if ok {

					if dna.CountBase(gV.Seq[gV.One.AlleleOne], dna.Gap) == int(read.Cigar[i].RunLen) && dna.CountBase(gV.Seq[gV.One.AlleleTwo], dna.Gap) == int(read.Cigar[i].RunLen) {
						parentAllele1++
					}
					if dna.CountBase(gV.Seq[gV.Two.AlleleOne], dna.Gap) == int(read.Cigar[i].RunLen) && dna.CountBase(gV.Seq[gV.Two.AlleleTwo], dna.Gap) == int(read.Cigar[i].RunLen) {
						parentAllele1++
					}
				}
				target += int(read.Cigar[i].RunLen)
			case 'M':
				for j = 0; j < int(read.Cigar[i].RunLen); j++ {
					code = vcf.GenomeKeyUint64(int(vcfHeader.Ref[read.RName]), int(target+j))
					gV, ok = snpDb[code]
					if ok {
						if string(read.Seq[query+j]) == string(gV.Seq[gV.One.AlleleOne]) && string(read.Seq[query+j]) == string(gV.Seq[gV.One.AlleleTwo]) {
							parentAllele1++
						}
						if string(read.Seq[query+j]) != string(gV.Seq[gV.Two.AlleleOne]) && string(read.Seq[query+j]) == string(gV.Seq[gV.Two.AlleleTwo]) {
							parentAllele2++
						}
					}

				}
				target += int(read.Cigar[i].RunLen)
				query += int(read.Cigar[i].RunLen)
			}
		}
		switch true {
		case parentAllele1 > parentAllele2:
			childOne.WriteString(bam.ToString(&read))
			childOne.WriteByte('\n')
		case parentAllele2 > parentAllele1:
			childTwo.WriteString(bam.ToString(&read))
			childTwo.WriteByte('\n')
		}
	}
	file.Reader.Close()
}
