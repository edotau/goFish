package genePred

import ()

func FilterStrand(gp []GenePred) []GenePred {
	var ans []GenePred
	for _, i := range gp {
		if i.Strand != '.' {
			ans = append(ans, i)
		}
	}
	return ans
}
