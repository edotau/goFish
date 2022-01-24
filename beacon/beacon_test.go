package beacon

import (
	"testing"
)

var testfile string = "testdata/D800000_OptoSeqBCRSummaryFile.csv"

var toyBeacon *OptoSeq = &OptoSeq{
	DeviceID:       "D82698",
	PenId:          5068,
	WellPlateID:    "WP30105",
	WellPlateIndex: 0,
	WellRow:        "A",
	WellColumn:     1,
	Barcode:        "C0D0F1T0",
	ExportCount:    1,
	Species:        "mm",
	BarcodeSet:     "a4_121820",
}

func TestBeacon(t *testing.T) {
	checkfile := Read(testfile)
	lines := 0

	opto := OptoSeq{}
	for i, val := range checkfile {
		t.Logf("%v\n", val)
		lines++
		if i == 0 {
			opto = val
		}
	}
	if lines != 546 {
		t.Fatalf("Error: Line numbers from beacon reader do not match %d != 546...\n", lines)
	}
	if opto.DeviceID != toyBeacon.DeviceID || opto.Barcode != toyBeacon.Barcode || opto.BarcodeSet != toyBeacon.BarcodeSet {
		t.Logf("Error: Device IDs do not match\n")
	}

}
