package main

import ()

const (
	chr01 int = 30751940
	chr02 int = 24095676
	chr03 int = 18730912
	chr04 int = 35952845
	chr05 int = 16512284
	chr06 int = 19184859
	chr07 int = 31778283
	chr08 int = 22385145
	chr09 int = 22564259
	chr10 int = 17640933
	chr11 int = 18132137
	chr12 int = 22068783
	chr13 int = 21534086
	chr14 int = 17482575
	chr15 int = 18148561
	chr16 int = 20569804
	chr17 int = 20701417
	chr18 int = 16681310
	chr19 int = 20543907
	chr20 int = 21246911
	chr21 int = 17891855
	chrM  int = 16713
)

func getChrom(name string) int {
	switch name {
	case "chr01":
		return chr01
	case "chr02":
		return chr02
	case "chr03":
		return chr03
	case "chr04":
		return chr04
	case "chr05":
		return chr05
	case "chr06":
		return chr06
	case "chr07":
		return chr07
	case "chr08":
		return chr08
	case "chr09":
		return chr09
	case "chr10":
		return chr10
	case "chr11":
		return chr11
	case "chr12":
		return chr12
	case "chr13":
		return chr13
	case "chr14":
		return chr14
	case "chr15":
		return chr15
	case "chr16":
		return chr16
	case "chr17":
		return chr17
	case "chr18":
		return chr18
	case "chr19":
		return chr19
	case "chr20":
		return chr20
	case "chr21":
		return chr21
	case "chrM":
		return chrM
	}
	return 0
}
