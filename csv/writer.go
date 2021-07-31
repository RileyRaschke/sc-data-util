package csv

import (
	"fmt"
	"io"
	"time"

	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util"
	log "github.com/sirupsen/logrus"
)

const CSV_HEADER = string("Date,Time,Open,High,Low,Last,Volume,NumTrades,BidVolume,AskVolume,PriorSettle")

type CsvBarRow struct {
	scid.IntradayRecord
	DateTime    time.Time
	PriorSettle float32
}

func (x CsvBarRow) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		x.DateTime.Format("2006/1/2"),
		x.DateTime.Format("15:04:05"),
		x.Open,
		x.High,
		x.Low,
		x.Close,
		x.NumTrades,
		x.TotalVolume,
		x.BidVolume,
		x.AskVolume,
		x.PriorSettle,
	)
}
func (x CsvBarRow) TickString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		x.DateTimeSC.Time().Format("2006/1/2"),
		x.DateTimeSC.Time().Format("15:04:05.000000"),
		x.Open,
		x.High,
		x.Low,
		x.Close,
		x.NumTrades,
		x.TotalVolume,
		x.BidVolume,
		x.AskVolume,
		x.PriorSettle,
	)
}

func DumpBarCsv(outFile interface{}, r *scid.ScidReader, startTime time.Time, endTime time.Time, barSize string) error {
	r.JumpTo(startTime)
	w, err := util.WriteBuffer(outFile)
	if err != nil {
		log.Errorf("Failed to open \"%v\" for writing with error: %v", outFile, err)
	}
	bDuration, err := time.ParseDuration(barSize)
	scdt_barStart := scid.NewSCDateTimeMs(startTime)
	scdt_endTime := scid.NewSCDateTimeMs(endTime)
	scdt_nextBar := scid.NewSCDateTimeMs(startTime.Add(bDuration))
	scdt_duration := scdt_nextBar - scdt_barStart
	scdt_nextBar = scdt_barStart // hacky, but efficient
	var barRow CsvBarRow
	w.WriteString(CSV_HEADER + "\n")
	for {
		rec, err := r.NextRecord()
		if err == io.EOF {
			if barRow.TotalVolume != 0 {
				w.WriteString(barRow.String() + "\n")
			}
			break
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
		if rec.DateTimeSC >= scdt_nextBar {
			if barRow.TotalVolume != 0 {
				w.WriteString(barRow.String() + "\n")
			}
			if rec.DateTimeSC >= scdt_endTime {
				break
			}
			scdt_barStart = scdt_nextBar
			for {
				if scdt_nextBar > rec.DateTimeSC {
					break
				} else {
					scdt_barStart = scdt_nextBar
					scdt_nextBar += scdt_duration
				}
			}
			barRow = CsvBarRow{IntradayRecord: *rec}
			barRow.DateTime = scdt_barStart.Time()
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
	w.Flush()
	return nil
}

func DumpRawTicks(outFile interface{}, r *scid.ScidReader, startTime time.Time, endTime time.Time, aggregation uint) {
	r.JumpTo(startTime)
	w, err := util.WriteBuffer(outFile)
	//scdt_startTime := scid.NewSCDateTimeMs(startTime)
	scdt_endTime := scid.NewSCDateTimeMs(endTime)
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
