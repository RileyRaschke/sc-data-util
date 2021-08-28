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
	scid.IntradayRecord
	DateTime    time.Time
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
func (x CsvBarRow) TickString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		x.DateTimeSC.Time().Format("2006/01/02"),
		x.DateTimeSC.Time().Format("15:04:05.000000"),
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

type BarAccumulator struct {
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	scdt_nextBar  scid.SCDateTimeMS
	scdt_duration scid.SCDateTimeMS
}

func NewBarAccumulator(startTime time.Time, endTime time.Time, barSize string) *BarAccumulator {
	x := BarAccumulator{}
	bDuration, _ := time.ParseDuration(barSize)
	x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
	x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
	x.scdt_nextBar = scid.NewSCDateTimeMS(startTime.Add(bDuration))
	x.scdt_duration = x.scdt_nextBar - x.scdt_barStart
	x.scdt_nextBar = x.scdt_barStart // hacky, but efficient
	return &x
}

func (x *BarAccumulator) AccumulateBar(r *scid.ScidReader) (CsvBarRow, error) {
	var barRow CsvBarRow
	for {
		rec, err := r.NextRecord()
		if err == io.EOF {
			//w.WriteString(barRow.String() + "\n")
			return barRow, err
		}
		if err != nil {
			log.Infof("Error returned by `r.NextRecord()`: %v", err)
		}
		if rec.Open == scid.SINGLE_TRADE_WITH_BID_ASK {

		} else if rec.Open == scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE {
		} else if rec.Open == scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE {
		} else {
			// support for index style data
			if rec.High == rec.Low {
				if rec.High < rec.Open {
					rec.High = rec.Open
				}
				if rec.High < rec.Close {
					rec.High = rec.Close
				}
				if rec.Low > rec.Open || (rec.Open != 0 && rec.Low < 0.5*rec.Open) {
					rec.Low = rec.Open
				}
				if rec.Low > rec.Close || (rec.Close != 0 && rec.Low < 0.5*rec.Close) {
					rec.Low = rec.Close
				}
			}
		}
		if rec.DateTimeSC >= x.scdt_nextBar {
			if barRow.TotalVolume != 0 {
				return barRow, nil
			}
			if rec.DateTimeSC >= x.scdt_endTime {
				break
			}
			x.scdt_barStart = x.scdt_nextBar
			for {
				if x.scdt_nextBar > rec.DateTimeSC {
					break
				} else {
					x.scdt_barStart = x.scdt_nextBar
					x.scdt_nextBar += x.scdt_duration
				}
			}
			barRow = CsvBarRow{IntradayRecord: *rec}
			barRow.DateTime = x.scdt_barStart.Time()
			barRow.Open = rec.Close
			//barRow.PriorSettle = getPriorSettle()? Nah.. Need to track in loop
		} else {
			if rec.High > barRow.High {
				barRow.High = rec.High
			}
			if rec.Low < barRow.Low {
				barRow.Low = rec.Low
			}
			barRow.Close = rec.Close
			barRow.NumTrades += rec.NumTrades
			barRow.TotalVolume += rec.TotalVolume
			barRow.BidVolume += rec.BidVolume
			barRow.AskVolume += rec.AskVolume
		}
	}
	return barRow, nil
}

func DumpBarCsv(outFile interface{}, r *scid.ScidReader, startTime time.Time, endTime time.Time, barSize string) error {
	r.JumpTo(startTime)
	w, err := util.WriteBuffer(outFile)
	if err != nil {
		log.Errorf("Failed to open \"%v\" for writing with error: %v", outFile, err)
	}
	w.WriteString(CSV_HEADER + "\n")
	ba := NewBarAccumulator(startTime, endTime, barSize)
	for {
		barRow, err := ba.AccumulateBar(r)
		if err != nil {
			if barRow.TotalVolume != 0 {
				w.WriteString(barRow.String() + "\n")
			}
			break
		}
		w.WriteString(barRow.String() + "\n")
	}
	w.Flush()
	return nil
}

func DumpRawTicks(outFile interface{}, r *scid.ScidReader, startTime time.Time, endTime time.Time, aggregation uint) {
	r.JumpTo(startTime)
	w, err := util.WriteBuffer(outFile)
	//scdt_startTime := scid.NewSCDateTimeMS(startTime)
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
		barRow := CsvBarRow{IntradayRecord: *rec}
		w.WriteString(fmt.Sprintf("%v\n", barRow.TickString()))
	}
	w.Flush()
}
