package vcf

/*
func GetAlleleGenotype(samples string) []Genotype {
	text := strings.Split(samples, "\t")
	var hap string
	var alleles []string
	var err error
	var n int64
	var answer []Genotype = make([]Genotype, len(text))
	for i := 0; i < len(text); i++ {
		hap = strings.Split(text[i], ":")[0]
		if strings.Compare(hap, "./.") == 0 || strings.Compare(hap, ".|.") == 0 {
			answer[i] = Genotype{One: Allele{Id: -1}, Two: Allele{Id: -1}, Phased: false}
		} else if strings.Contains(hap, "|") {
			alleles = strings.SplitN(hap, "|", 2)
			answer[i] = Genotype{One: Allele{Id: int16(simpleio.StringToInt(alleles[0]))}, Two: Allele{Id: int16(simpleio.StringToInt(alleles[1]))}, Phased: true}
		} else if strings.Contains(hap, "/") {
			alleles = strings.SplitN(hap, "/", 2)
			answer[i] = Genotype{One: Allele{Id: int16(simpleio.StringToInt(alleles[0]))}, Two: Allele{Id: int16(simpleio.StringToInt(alleles[1]))}, Phased: false}
		} else {
			//Deal with single haps. There might be a better soltuion, but I think this should work.
			n, err = strconv.ParseInt(alleles[0], 10, 16)
			if err != nil && n < int64(len(text)) {
				answer[i] = Genotype{One: Allele{Id: int16(n)}, Two: Allele{Id: -1}, Phased: false}
			} else {
				log.Fatalf("Error: Unexpected parsing error...\n")
			}
		}
	}
	return answer
}*/
