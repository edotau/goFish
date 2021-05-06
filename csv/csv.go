// Package csv contains data structures to parse comma separated value takes
// and perform statical analysis
package csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/edotau/goFish/simpleio"
)

func ReadToXY(filename string) ([]float64, []float64) {
	reader := simpleio.NewReader(filename)
	var x, y []float64
	var line []string
	for i, err := simpleio.ReadLine(reader); !err; i, err = simpleio.ReadLine(reader) {
		line = strings.Split(i.String(), ",")
		x = append(x, simpleio.StringToFloat64(line[0]))
		y = append(y, simpleio.StringToFloat64(line[1]))
	}
	return x, y
}

func Write(filepath string, x [][]float64, y []float64, highPrecision bool) {
	lenX := len(x)
	lenY := len(y)
	var numFeatures int

	if lenX != 0 {
		numFeatures = len(x[0])
	}

	if lenX == 0 || lenY == 0 || numFeatures == 0 || lenX != lenY {

		log.Fatal(fmt.Errorf("ERROR: Training set (either x or y or both) has no examples or the lengths of the dataset don't match\n\tlength of x: %v\n\tlength of y: %v\n\tnumber of features in x: %v\n", lenX, lenY, numFeatures))
	}

	_, err := os.Stat(filepath)
	if err != nil && !os.IsNotExist(err) {
		simpleio.StdError(err)
	}

	file, err := os.Create(filepath)
	if err != nil {
		simpleio.StdError(err)
	}
	defer file.Close()

	var precision int
	if highPrecision {
		precision = 64
	} else {
		precision = 32
	}

	writer := csv.NewWriter(file)
	records := [][]string{}

	// parse until the end of the file
	for i := range x {
		record := []string{}

		for j := range x[i] {
			record = append(record, strconv.FormatFloat(x[i][j], 'g', -1, precision))
		}

		record = append(record, strconv.FormatFloat(y[i], 'g', -1, precision))

		records = append(records, record)
	}
	// now save the record to file
	err = writer.WriteAll(records)
	simpleio.StdError(err)
}

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
