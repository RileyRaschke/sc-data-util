package csv

import (
	"fmt"
	"io"
	"time"

	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util"
	log "github.com/sirupsen/logrus"
)

const CSV_HEADER = string("Date,Time,Open,High,Low,Last,Volume,NumTrades,BidVolume,AskVolume")
const CSV_HEADER_DETAIL = string("Date,Time,Open,High,Low,Last,Volume,NumTrades,BidVolume,AskVolume,PriorLast,PriorSettle,TradingDate")

type CsvBarRow struct {
	util.BasicBar
	PriorSettle float32
	PriorLast   float32
	TradingDate time.Time
}

func (x CsvBarRow) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		x.DateTime.Format("2006/01/02"),
		x.DateTime.Format("15:04:05"),
		x.Open,
		x.High,
		x.Low,
		x.Close,
		x.NumTrades,
		x.TotalVolume,
		x.BidVolume,
		x.AskVolume,
	)
}
func (x CsvBarRow) DetailString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		x.DateTime.Format("2006/01/02"),
		x.DateTime.Format("15:04:05"),
		x.Open,
		x.High,
		x.Low,
		x.Close,
		x.NumTrades,
		x.TotalVolume,
		x.BidVolume,
		x.AskVolume,
		x.PriorLast,
		x.PriorSettle,
		x.TradingDate.Format("2006/01/02"),
	)
}

func DumpBarCsv(outFile interface{}, r *scid.ScidReader, startTime time.Time, endTime time.Time, barSize string) error {
	r.JumpTo(startTime)
	w, err := util.WriteBuffer(outFile)
	if err != nil {
		log.Errorf("Failed to open \"%v\" for writing with error: %v", outFile, err)
	}
	w.WriteString(CSV_HEADER + "\n")
	ba := util.NewBarAccumulator(startTime, endTime, barSize)
	for {
		bar, err := ba.AccumulateBar(r)
		barRow := CsvBarRow{BasicBar: bar.(util.BasicBar)}
		if barRow.TotalVolume != 0 {
			w.WriteString(barRow.String() + "\n")
		}
		if err != nil {
			break
		}
	}
	w.Flush()
	return nil
}

func DumpRawTicks(outFile interface{}, r *scid.ScidReader, startTime time.Time, endTime time.Time, aggregation uint) {
	r.JumpTo(startTime)
	w, err := util.WriteBuffer(outFile)
	scdt_endTime := scid.NewSCDateTimeMS(endTime)
	if err != nil {
		log.Errorf("Failed to open \"%v\" for writing with error: %v", outFile, err)
	}
	w.WriteString(CSV_HEADER + "\n")
	for {
		rec, err := r.NextRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Infof("Error returned by `r.NextRecord()`: %v", err)
		}
		if rec.DateTimeSC >= scdt_endTime {
			break
		}
		barRow := CsvBarRow{BasicBar: util.BasicBar{IntradayRecord: *rec}}
		w.WriteString(fmt.Sprintf("%v\n", barRow.TickString()))
	}
	w.Flush()
}
