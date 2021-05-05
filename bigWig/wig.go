package bigWig

// Wig contains chromosome location and step properties of Wig data. Individual wig values are stored in the underlying WigValue struct
type Wig struct {
	StepType string    `bin:"len:2"`
	Chrom    string    `bin:"len:4"`
	Start    int       `bin:"len:4"`
	Step     int       `bin:"len:4"`
	Span     int       `bin:"len:4"`
	Val      []float64 `bin:"len:4"`
}

type WigValue struct {
	Position int
	Value    float64
}
