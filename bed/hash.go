package bed

import (
	"log"
	"strings"

	"github.com/edotau/goFish/reference"
	"github.com/edotau/goFish/simpleio"
)

func GenomeKey(div int) map[string][]*GenomeInfo {
	var region int
	genome := make(map[string][]*GenomeInfo)
	for _, chrom := range reference.Chr {
		region = reference.GetChrom(chrom) / div
		genome = divide(chrom, region, genome)
	}
	return genome
}

func buildreference(coord *GenomeInfo, keys map[string][]*GenomeInfo, genome map[string][]*GenomeInfo) map[string][]*GenomeInfo {
	for _, region := range keys[coord.Chr] {
		if Overlap(region, coord) {

			genome[region.Info.String()] = append(genome[region.Info.String()], coord)
			return genome
		}
	}
	return genome
}

func divide(chrom string, size int, chromHash map[string][]*GenomeInfo) map[string][]*GenomeInfo {
	var i int
	for i = 0; i < reference.GetChrom(chrom)-size; i += size {
		curr := GenomeInfo{Chr: chrom, Start: i, End: i + size}
		curr.Info.WriteString(curr.Chrom() + simpleio.IntToString(curr.ChrEnd()))

		chromHash[chrom] = append(chromHash[chrom], &curr)
	}
	last := GenomeInfo{Chr: chrom, Start: i, End: reference.GetChrom(chrom)}
	last.Info.WriteString(last.Chrom() + simpleio.IntToString(last.ChrEnd()))
	chromHash[chrom] = append(chromHash[chrom], &last)
	return chromHash
}

func SelectGenomeHash(filename string, num int) (map[string][]*GenomeInfo, map[string][]*GenomeInfo) {
	keys := GenomeKey(num)
	answer := make(map[string][]*GenomeInfo)
	selectRegions := simpleio.NewReader(filename)
	for i, err := ToGenomeInfo(selectRegions); !err; i, err = ToGenomeInfo(selectRegions) {
		answer = buildreference(i, keys, answer)
	}
	return keys, answer
}

func CheckOverlapHash(region string, b *GenomeInfo, answer map[string][]*GenomeInfo) bool {
	for _, r := range answer[region] {
		if Overlap(r, b) {
			return true
		}
	}
	return false
}

func OverlapHashFilterSv(region string, b *GenomeInfo, answer map[string][]*GenomeInfo, filterSv string) bool {
	for _, r := range answer[region] {
		if Overlap(r, b) && strings.Contains(r.Info.String(), filterSv) {
			return true
		}
	}
	return false
}

func GetHashKey(b *GenomeInfo, keys map[string][]*GenomeInfo) string {
	for _, each := range keys[b.Chrom()] {
		if Overlap(b, each) {
			return each.Info.String()
		}
	}
	log.Fatalf("Error: query needs to match at least one select region...\n")
	return ""
}
