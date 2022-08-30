package util

import (
	"fmt"
	"io"
	"time"

	"github.com/RileyR387/sc-data-util/scid"
	log "github.com/sirupsen/logrus"
)

type Profile interface {
	String() string
}

type BasicProfileBar struct {
	Price       float32
	BidVolume   uint32
	AskVolume   uint32
	TotalVolume uint32
	BidTrades   uint32
	AskTrades   uint32
	NumTrades   uint32
}

type BasicProfile struct {
	DateTime       time.Time
	Duration       time.Time
	ProfileDataMap map[float32]*BasicProfileBar
}

func (x *BasicProfile) String() string {
	return "FIXME"
}

func (x *BasicProfileBar) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v",
		x.Price,
		x.BidVolume,
		x.AskVolume,
		x.TotalVolume,
		x.BidTrades,
		x.AskTrades,
		x.NumTrades,
	)
}
func (x *BasicProfileBar) TickString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v",
		x.Price,
		x.BidVolume,
		x.AskVolume,
		x.TotalVolume,
		x.BidTrades,
		x.AskTrades,
		x.NumTrades,
	)
}

type ProfileAccumulator struct {
	*BasicProfile
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	sctd_nextBar  scid.SCDateTimeMS
	scdt_duration scid.SCDateTimeMS
}

func NewProfileAccumulator(startTime time.Time, endTime time.Time, barSize string) *ProfileAccumulator {
	x := ProfileAccumulator{}
	bDuration, _ := time.ParseDuration(barSize)
	x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
	x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
	x.sctd_nextBar = scid.NewSCDateTimeMS(startTime.Add(bDuration))
	x.scdt_duration = x.sctd_nextBar - x.scdt_barStart
	x.sctd_nextBar = x.scdt_barStart // hacky, but efficient
	return &x
}

func (x *ProfileAccumulator) AccumulateTick(r *scid.ScidReader) (Profile, error) {
	var barRow BasicProfileBar
	x.BasicProfile = &BasicProfile{}
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return x.BasicProfile, nil
		}
		if rec.Open == scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			log.Info("scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE - Unhandled")
			log.Info(rec)
			continue
		} else if rec.Open == scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			log.Info("scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE - Unhandled")
			log.Info(rec)
			continue
		} else if rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
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
				return x.BasicProfile, nil
			}
			if rec.DateTimeSC >= x.scdt_endTime {
				return x.BasicProfile, io.EOF
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
			//barRow = BasicProfileBar{IntradayRecord: *rec}
			//barRow.DateTime = x.scdt_barStart.Time()
		} else {
			barRow.Price = rec.Close
			if rec.High == rec.Low {
			} else {
			}
			barRow.NumTrades += rec.NumTrades
			barRow.TotalVolume += rec.TotalVolume
			barRow.BidVolume += rec.BidVolume
			barRow.AskVolume += rec.AskVolume
		}
	}
	return x.BasicProfile, nil
}
