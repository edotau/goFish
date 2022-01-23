package beacon

import (
	"log"
	"strconv"
	"strings"

	"github.com/edotau/goFish/simpleio"
)

type OptoSeq struct {
	DeviceID       string
	PenId          int
	WellPlateID    string
	WellPlateIndex int
	WellRow        string
	WellColumn     int
	Barcode        string
	ExportCount    int
	Species        string
	BarcodeSet     string
}

func (b *OptoSeq) String() string {
	str := strings.Builder{}
	str.WriteString(b.DeviceID)
	str.WriteByte(',')
	str.WriteString(strconv.Itoa(b.PenId))
	str.WriteByte(',')
	str.WriteString(b.WellPlateID)
	str.WriteByte(',')
	str.WriteString(strconv.Itoa(b.WellPlateIndex))
	str.WriteString(b.WellRow)
	str.WriteByte(',')
	str.WriteString(strconv.Itoa(b.WellColumn))
	str.WriteByte(',')
	str.WriteString(b.Barcode)
	str.WriteByte(',')
	str.WriteString(strconv.Itoa(b.ExportCount))
	str.WriteByte(',')
	str.WriteString(b.Species)
	str.WriteByte(',')
	str.WriteString(b.BarcodeSet)

	return str.String()
}

func Read(filename string) []OptoSeq {
	var ans []OptoSeq
	//var err bool
	reader := simpleio.NewReader(filename)
	reader.Buffer, _ = simpleio.ReadLine(reader)
	//if !err && reader.Buffer.String() != "DeviceID,PenId,WellPlateID,WellPlateIndex,WellRow,WellColumn,Barcode,ExportCount,Species,BarcodeSet" {
	//	log.Fatal("Error: Beacon file is missing header line...\n")
	//}
	for i, done := parseBeaconData(reader); !done; i, done = parseBeaconData(reader) {
		ans = append(ans, *i)
	}
	return ans
}

func parseBeaconData(reader *simpleio.SimpleReader) (*OptoSeq, bool) {
	var err bool
	reader.Buffer, err = simpleio.ReadLine(reader)

	if !err {
		columns := strings.Split(reader.Buffer.String(), ",")
		if len(columns) != 10 {
			log.Fatal("Error: Beacon file contains exactly 10 fields...\n")
		}
		return &OptoSeq{
			DeviceID:       columns[0],
			PenId:          simpleio.StringToInt(columns[1]),
			WellPlateID:    columns[2],
			WellPlateIndex: simpleio.StringToInt(columns[3]),
			WellRow:        columns[4],
			WellColumn:     simpleio.StringToInt(columns[5]),
			Barcode:        columns[6],
			ExportCount:    simpleio.StringToInt(columns[7]),
			Species:        columns[8],
			BarcodeSet:     columns[9],
		}, false
	} else {
		return nil, true
	}
}

func replaceRef(record *OptoSeq) *OptoSeq {
	record.Species = "dog"
	return record
}
