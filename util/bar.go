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

func NewBasicBar(rec *scid.IntradayRecord, scdt scid.SCDateTimeMS) BasicBar {
	x := BasicBar{scid.IntradayRecord{}, scdt.Time()}
	x.High = rec.Close
	x.Low = rec.Close
	x.DateTimeSC = scdt
	x.Open = rec.Close
	x.Close = rec.Close
	x.NumTrades += rec.NumTrades
	x.TotalVolume += rec.TotalVolume
	x.BidVolume += rec.BidVolume
	x.AskVolume += rec.AskVolume
	return x
}

func (x *BasicBar) AddRecord(rec *scid.IntradayRecord) {
	if rec.Close > x.High {
		x.High = rec.Close
	}
	if rec.Close < x.Low {
		x.Low = rec.Close
	}
	x.Close = rec.Close
	x.NumTrades += rec.NumTrades
	x.TotalVolume += rec.TotalVolume
	x.BidVolume += rec.BidVolume
	x.AskVolume += rec.AskVolume
}
