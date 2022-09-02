package util

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util/bartype"
	log "github.com/sirupsen/logrus"
)

type BarAccumulator interface {
	AccumulateBar(*scid.ScidReader) (Bar, error)
}

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

type TimeBarAccumulator struct {
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	scdt_nextBar  scid.SCDateTimeMS
	scdt_duration scid.SCDateTimeMS
	barType       bartype.BarType
	barSize       int64
}

type VolumeBarAccumulator struct {
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	barSize       int64
	remainder     BasicBar
}
type TickBarAccumulator struct {
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	barSize       int64
	remainder     BasicBar
}

func NewBarAccumulator(startTime time.Time, endTime time.Time, barSize string) BarAccumulator {
	bt, duration := parseBarSize(barSize)
	switch bt {
	case bartype.Time:
		x := TimeBarAccumulator{}
		x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
		x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
		x.scdt_nextBar = scid.NewSCDateTimeMS(startTime.Add(time.Duration(duration)))
		x.scdt_duration = x.scdt_nextBar - x.scdt_barStart
		x.scdt_nextBar = x.scdt_barStart // hacky, but efficient
		return &x
	case bartype.Tick:
		x := TickBarAccumulator{}
		x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
		x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
		return &x
	case bartype.Volume:
		x := VolumeBarAccumulator{}
		x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
		x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
		return &x
	}
	return &TimeBarAccumulator{}
}

func parseBarSize(barSize string) (bartype.BarType, int64) {
	t := bartype.ParseType(barSize)
	var duration int64
	if t != bartype.Time {
		duration, _ = strconv.ParseInt(barSize[0:len(barSize)-1], 10, 64)
	} else {
		d, _ := time.ParseDuration(barSize)
		duration = int64(d)
	}
	return t, duration
}

func (x *TickBarAccumulator) AccumulateBar(r *scid.ScidReader) (Bar, error) {
	var barRow BasicBar
	var unbundled = false
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return barRow, err
		}
		if rec.Open == scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			//err := bundleTrades(r, rec)

			//log.Trace("FIXME?: scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE - Unhandled")
			//log.Trace(rec)
			log.Tracef("First Sub: %s", rec)
			unbundled = true
			continue
		} else if rec.Open == scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			//log.Trace("FIXME?: scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE - Unhandled")
			//log.Trace(rec)
			log.Tracef("Last Sub: %s", rec)
			unbundled = false
			continue
		} else if rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
			normalizeIndexData(rec)
		} else if unbundled {
			log.Tracef("between: %s", rec)
		} else {
			log.Tracef("trade: %s", rec)
		}
		if rec.DateTimeSC >= x.scdt_endTime {
			return barRow, io.EOF
		}
	}
	return barRow, nil
}

func (x *VolumeBarAccumulator) AccumulateBar(r *scid.ScidReader) (Bar, error) {
	var barRow BasicBar
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return barRow, err
		}
		if rec.Open == scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			err := bundleTrades(r, rec)
			if err != nil {
				return barRow, err
			}
			continue
		} else if rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
			normalizeIndexData(rec)
		}
		if rec.DateTimeSC >= x.scdt_endTime {
			return barRow, io.EOF
		}
	}
	return barRow, nil
}

func (x *TimeBarAccumulator) AccumulateBar(r *scid.ScidReader) (Bar, error) {
	var barRow BasicBar
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return barRow, err
		}
		if rec.Open == scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			log.Info("FIXME?: scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE - Unhandled")
			log.Info(rec)
			continue
		} else if rec.Open == scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			log.Info("FIXME?: scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE - Unhandled")
			log.Info(rec)
			continue
		} else if rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
			normalizeIndexData(rec)
		}

		if rec.DateTimeSC >= x.scdt_nextBar {
			if barRow.TotalVolume != 0 {
				return barRow, nil
			}
			if rec.DateTimeSC >= x.scdt_endTime {
				return barRow, io.EOF
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

func bundleTrades(r *scid.ScidReader, bundle *scid.IntradayRecord) error {
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return err
		}
		bundle.TotalVolume += rec.TotalVolume
		bundle.BidVolume += rec.BidVolume
		bundle.AskVolume += rec.AskVolume
		bundle.NumTrades += rec.NumTrades
		if rec.High > bundle.High {
			bundle.High = rec.High
		}
		if rec.Low < bundle.Low {
			bundle.Low = rec.Low
		}
		if rec.Open == scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			// assume the last record is the correct close
			bundle.Close = rec.Close
			log.Tracef("Bundled trade: %s", rec)
			return nil
		}
	}
	return nil
}

func normalizeIndexData(rec *scid.IntradayRecord) {
	// support for index style data
	if rec.High == rec.Low {
		if rec.High < rec.Open {
			log.Debugf("High(%f) is below the Open(%f) at %s", rec.High, rec.Open, rec.DateTimeSC)
			rec.High = rec.Open
		}
		if rec.High < rec.Close {
			log.Debugf("High(%f) is below the Close(%f) at %s", rec.High, rec.Open, rec.DateTimeSC)
			rec.High = rec.Close
		}
		if rec.Low > rec.Open {
			log.Debugf("Low(%f) is above the Open(%f) at %s", rec.Low, rec.Open, rec.DateTimeSC)
			rec.Low = rec.Open
		}
		if rec.Open != 0 && rec.Low < 0.95*rec.Open {
			log.Debugf("Low(%f) is 95%% below Open(%f) at %s", rec.Low, rec.Open, rec.DateTimeSC)
			rec.Low = rec.Open
		}
		if rec.Low > rec.Close {
			log.Debugf("Low(%f) is above the Close(%f) at %s", rec.Low, rec.Close, rec.DateTimeSC)
			rec.Low = rec.Close
		}
		if rec.Close != 0 && rec.Low < 0.95*rec.Close {
			log.Debugf("Low(%f) is 95%% below Close(%f) at %s", rec.Low, rec.Close, rec.DateTimeSC)
			rec.Low = rec.Close
		}
	}
}
