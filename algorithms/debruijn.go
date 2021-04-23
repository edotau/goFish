// Package algorithms contains basic implementations of popular algorithms in computer science
package algorithms

// DeDruijn is a function returns a sequence for a n-words with k letters, modifiied from the Biogo Authors, copyright ©2011-2012
func DeBruijn(k, n byte) (s []byte) {
	switch k {
	case 0:
		return []byte{}
	case 1:
		return make([]byte, n)
	}

	a := make([]byte, k*n)
	s = make([]byte, 0, Pow(int(k), n))

	var db func(byte, byte)
	db = func(t, p byte) {
		if t > n {
			if n%p == 0 {
				for j := byte(1); j <= p; j++ {
					s = append(s, a[j])
				}
			}
		} else {
			a[t] = a[t-p]
			db(t+1, p)
			for j := a[t-p] + 1; j < k; j++ {
				a[t] = j
				db(t+1, t)
			}

		}
	}
	db(1, 1)

	return
}

// Return the exp'th power of base.
func Pow(base int, exp byte) (r int) {
	r = 1
	for exp > 0 {
		if exp&1 != 0 {
			r *= base
		}
		exp >>= 1
		base *= base
	}

	return
}
