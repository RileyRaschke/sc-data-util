package util

import (
	"fmt"
	"io"
	"time"

	"github.com/RileyR387/sc-data-util/scid"
)

type Bar interface {
	String() string
	TickString() string
}

type BasicBar struct {
	scid.IntradayRecord
	DateTime time.Time
}

func (x BasicBar) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v",
		x.DateTime,
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
func (x BasicBar) TickString() string {
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
	sctd_nextBar  scid.SCDateTimeMS
	scdt_duration scid.SCDateTimeMS
}

func NewBarAccumulator(startTime time.Time, endTime time.Time, barSize string) *BarAccumulator {
	x := BarAccumulator{}
	bDuration, _ := time.ParseDuration(barSize)
	x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
	x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
	x.sctd_nextBar = scid.NewSCDateTimeMS(startTime.Add(bDuration))
	x.scdt_duration = x.sctd_nextBar - x.scdt_barStart
	x.sctd_nextBar = x.scdt_barStart // hacky, but efficient
	return &x
}

func (x *BarAccumulator) AccumulateBar(r *scid.ScidReader) (Bar, error) {
	var barRow BasicBar
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return barRow, err
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
		if rec.DateTimeSC >= x.sctd_nextBar {
			if barRow.TotalVolume != 0 {
				return barRow, nil
			}
			if rec.DateTimeSC >= x.scdt_endTime {
				return barRow, io.EOF
			}
			x.scdt_barStart = x.sctd_nextBar
			for {
				if x.sctd_nextBar > rec.DateTimeSC {
					break
				} else {
					x.scdt_barStart = x.sctd_nextBar
					x.sctd_nextBar += x.scdt_duration
				}
			}
			barRow = BasicBar{IntradayRecord: *rec}
			barRow.DateTime = x.scdt_barStart.Time()
			barRow.Open = rec.Close
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
