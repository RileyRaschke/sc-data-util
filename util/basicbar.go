package util

import (
	"fmt"
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
		x.TotalVolume,
		x.NumTrades,
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
		x.TotalVolume,
		x.NumTrades,
		x.BidVolume,
		x.AskVolume,
	)
}

func updateBar(barRow *BasicBar, rec *scid.IntradayRecord) {
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
