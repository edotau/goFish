package axt

import (
	"github.com/edotau/goFish/bed"
	"github.com/edotau/goFish/code"
	"github.com/edotau/goFish/simpleio"
)

func AxtToGenomeInfo(axtFile *Axt) []bed.GenomeInfo {
	var answer []bed.GenomeInfo

	var rCount int = axtFile.RStart - 1
	qCount := axtFile.QStart - 1
	for i := 0; i < len(axtFile.RSeq); i++ {
		if axtFile.RSeq[i] != '-' && axtFile.QSeq[i] != '-' {
			rCount++
			qCount++
			continue
			//snp mismatch

		}
		//insertion in VCF record
		if axtFile.RSeq[i] == '-' {

			qCount++

			target := bed.GenomeInfo{Chr: axtFile.RName, Start: rCount, End: rCount}

			query := bed.GenomeInfo{Chr: axtFile.QName, Start: qCount, End: qCount}
			target.Info.WriteString("INS")
			target.Info.WriteByte('\t')

			//curr.Info.WriteString(axtFile.QName)
			//curr.Info.WriteByte('\t')
			//curr.Info.WriteString(simpleio.IntToString(qCount))
			//curr.Info.WriteByte('\t')
			for j := i; j < len(axtFile.RSeq); j++ {
				if axtFile.RSeq[j] == code.Gap {
					query.End++
					qCount++
				} else {
					//curr = &vcf.Vcf{Chr: axtFile.RName, Pos: rCount, Id: axtFile.QName, Ref: dna.BaseToString(dna.ToUpper(axtFile.RSeq[i-1])), Alt: altTmp, Qual: 0, Filter: "PASS", Info: infoTag, Format: "SVTYPE=INS", Unknown: "GT:DP:AD:RO:QR:AO:QA:GL"}
					//query.End = qCount
					target.Info.WriteString(bed.GenomeInfoToString(query))
					target.Info.WriteByte('\t')
					target.Info.WriteByte(axtFile.QStrandPos)
					target.Info.WriteByte('\t')
					if diff := query.End - query.Start; diff > 10 {
						target.Info.WriteString(simpleio.IntToString(diff))
						answer = append(answer, target)
					}
					i = j - 1
					break
				}
			}
		}
		//deleteion vcf record
		if axtFile.QSeq[i] == '-' {
			tempRCount := 0
			target := bed.GenomeInfo{Chr: axtFile.RName, Start: rCount, End: rCount}
			query := bed.GenomeInfo{Chr: axtFile.QName, Start: qCount, End: qCount}
			target.Info.WriteString("DEL")
			target.Info.WriteByte('\t')

			//altTmp = dna.BaseToString(dna.ToUpper(axtFile.RSeq[i-1]))
			for j := i; j < len(axtFile.RSeq); j++ {
				if axtFile.QSeq[j] == code.Gap {

					//altTmp = altTmp +
					tempRCount++
				} else {

					rCount = rCount + tempRCount

					target.End = rCount
					target.End += tempRCount

					target.Info.WriteString(bed.GenomeInfoToString(query))
					target.Info.WriteByte('\t')
					target.Info.WriteByte(axtFile.QStrandPos)
					target.Info.WriteByte('\t')
					if diff := target.End - target.Start; diff > 10 {
						target.Info.WriteString(simpleio.IntToString(diff))
						answer = append(answer, target)
					}

					i = j - 1
					break
				}
			}
		}
	}
	//log.Printf("\nFound %d differences in this block...\n%s\n", len(answer), AxtInfo(axtFile))
	return answer
}
