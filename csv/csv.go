package csv

import ()

type FisherExact struct {
	A int
	B int
	C int
	D int
}

type PeakStats struct {
	Chr         string
	Start       int
	End         int
	Score       int
	Matrix      *FisherExact
	LeftPvalue  float64
	RightPvalue float64
}

type DiffPeak struct {
	Chr   string
	Start int
	End   int
	Pval  float64
}
