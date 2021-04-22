package vcf

import ()

func ASFilter(v *Vcf, parentOne int, parentTwo int, F1 int) bool {
	gtOne := ParseGt(v.Genotypes[parentOne])
	gtTwo := ParseGt(v.Genotypes[parentTwo])
	f1Hybrid := ParseGt(v.Genotypes[F1])
	if IsHomozygous(gtOne) && IsHomozygous(gtTwo) && IsHeterozygous(f1Hybrid) && gtOne.AlleleOne != gtTwo.AlleleOne {
		return true
	} else {
		return false
	}
}

func IsHeterozygous(genome Genotype) bool {

	if genome.AlleleOne < 0 || genome.AlleleTwo < 0 {
		return false
	} else if genome.AlleleOne != genome.AlleleTwo {
		return true
	} else if genome.AlleleOne == genome.AlleleTwo {
		return false
	} else {
		return false
	}
}

func IsHomozygous(genome Genotype) bool {
	if genome.AlleleOne < 0 || genome.AlleleTwo < 0 {
		return false
	}
	if genome.AlleleOne == genome.AlleleTwo {
		return true
	}
	if genome.AlleleOne != genome.AlleleTwo {
		return false
	}
	return false
}
