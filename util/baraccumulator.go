package util

import (
	"fmt"
	"io"
	"time"

	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util/bartype"
	log "github.com/sirupsen/logrus"
)

type BarAccumulator interface {
	AccumulateBar(*scid.ScidReader) (Bar, error)
}
type BarProfileAccumulator interface {
	AccumulateProfile(*scid.ScidReader) (Bar, BarProfile, error)
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
	bundle        bool
	withProfile   bool
	nextBar       BasicBar
}

type TickBarAccumulator struct {
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	barSize       uint32
	bundle        bool
	withProfile   bool
	nextBar       BasicBar
}

type VolumeBarAccumulator struct {
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	barSize       uint32
	withProfile   bool
	nextBar       BasicBar
}

func NewBarAccumulator(startTime time.Time, endTime time.Time, barSize string, bundleOpt bool, withProfile bool) BarAccumulator {
	bt, duration := parseBarSize(barSize)
	switch bt {
	case bartype.Time:
		x := TimeBarAccumulator{}
		x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
		x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
		x.scdt_nextBar = scid.NewSCDateTimeMS(startTime.Add(time.Duration(duration)))
		x.scdt_duration = x.scdt_nextBar - x.scdt_barStart
		x.scdt_nextBar = x.scdt_barStart // hacky, but efficient
		x.bundle = bundleOpt
		x.withProfile = withProfile
		return &x
	case bartype.Tick:
		x := TickBarAccumulator{}
		x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
		x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
		x.bundle = bundleOpt
		x.withProfile = withProfile
		x.barSize = uint32(duration) // in ticks
		return &x
	case bartype.Volume:
		x := VolumeBarAccumulator{}
		x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
		x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
		x.withProfile = withProfile
		x.barSize = uint32(duration) // in volume
		return &x
	}
	return &TimeBarAccumulator{}
}

func NewBarProfileAccumulator(startTime time.Time, endTime time.Time, barSize string, bundleOpt bool, withProfile bool) BarProfileAccumulator {
	bt, duration := parseBarSize(barSize)
	switch bt {
	case bartype.Time:
		x := TimeBarAccumulator{}
		x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
		x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
		x.scdt_nextBar = scid.NewSCDateTimeMS(startTime.Add(time.Duration(duration)))
		x.scdt_duration = x.scdt_nextBar - x.scdt_barStart
		x.scdt_nextBar = x.scdt_barStart // hacky, but efficient
		x.bundle = bundleOpt
		x.withProfile = withProfile
		return &x
	case bartype.Tick:
		x := TickBarAccumulator{}
		x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
		x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
		x.bundle = bundleOpt
		x.withProfile = withProfile
		x.barSize = uint32(duration) // in ticks
		return &x
	case bartype.Volume:
		x := VolumeBarAccumulator{}
		x.scdt_barStart = scid.NewSCDateTimeMS(startTime)
		x.scdt_endTime = scid.NewSCDateTimeMS(endTime)
		x.withProfile = withProfile
		x.barSize = uint32(duration) // in volume
		return &x
	}
	return &TimeBarAccumulator{}
}

// Time bars should typically bundle..
func (x *TimeBarAccumulator) AccumulateBar(r *scid.ScidReader) (Bar, error) {
	var barRow BasicBar
	if x.nextBar.TotalVolume > 0 {
		barRow = x.nextBar
		x.nextBar = BasicBar{}
	}
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return barRow, err
		}
		if x.bundle && rec.Open == scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			err := bundleTrades(r, rec)
			if err != nil {
				log.Warn("Error occured before trade was bundled!")
				return barRow, err
			}
		} else if rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
			normalizeIndexData(rec)
		}

		if rec.DateTimeSC >= x.scdt_endTime {
			return barRow, io.EOF
		}

		if rec.DateTimeSC >= x.scdt_nextBar {
			x.scdt_barStart = x.scdt_nextBar
			// assure the next tick is within the next bar's duration
			for {
				if x.scdt_nextBar > rec.DateTimeSC {
					break
				} else {
					x.scdt_barStart = x.scdt_nextBar
					x.scdt_nextBar += x.scdt_duration
				}
			}
			x.nextBar = BasicBar{IntradayRecord: *rec}
			x.nextBar.DateTime = x.scdt_barStart.Time()
			x.nextBar.Open = rec.Close

			if barRow.TotalVolume != 0 {
				return barRow, nil
			}
		} else {
			updateBar(&barRow, rec)
		}
	}
	return barRow, nil
}

func (x *TimeBarAccumulator) AccumulateProfile(r *scid.ScidReader) (Bar, BarProfile, error) {
	var barRow BasicBar
	var barProfile BarProfile
	return barRow, barProfile, nil
}

func (x *TickBarAccumulator) AccumulateProfile(r *scid.ScidReader) (Bar, BarProfile, error) {
	var barRow BasicBar
	var barProfile BarProfile
	return barRow, barProfile, nil
}
func (x *VolumeBarAccumulator) AccumulateProfile(r *scid.ScidReader) (Bar, BarProfile, error) {
	var barRow BasicBar
	var barProfile BarProfile
	return barRow, barProfile, nil
}

// Tick bars should support optional bundleing
func (x *TickBarAccumulator) AccumulateBar(r *scid.ScidReader) (Bar, error) {
	var barRow BasicBar

	rec, err := r.NextRecord()
	if err != nil {
		return barRow, err
	}

	if x.nextBar.TotalVolume > 0 {
		barRow = x.nextBar
		x.nextBar = BasicBar{}
	} else {
		barRow = BasicBar{IntradayRecord: *rec}
		barRow.DateTime = rec.DateTimeSC.Time()
		barRow.Open = rec.Close
	}

	for {
		rec, err := r.NextRecord()
		if err != nil {
			return barRow, err
		}

		if x.bundle && rec.Open == scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			err := bundleTrades(r, rec)
			if err != nil {
				log.Warn("Error occured before trade was bundled!")
				return barRow, err
			}
		} else if rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
			normalizeIndexData(rec)
		}

		if rec.DateTimeSC >= x.scdt_endTime {
			return barRow, io.EOF
		}

		updateBar(&barRow, rec)

		if barRow.NumTrades < x.barSize {
			continue
		}
		if barRow.NumTrades == x.barSize {
			x.nextBar = BasicBar{}
			return barRow, nil
		}
		if barRow.NumTrades > x.barSize {
			overage := barRow.NumTrades - x.barSize

			x.nextBar = BasicBar{IntradayRecord: *rec}
			x.nextBar.DateTime = rec.DateTimeSC.Time()
			x.nextBar.Open = rec.Close

			barRow.TotalVolume -= overage
			x.nextBar.TotalVolume = overage

			if rec.BidVolume >= rec.AskVolume && rec.BidVolume >= overage {
				barRow.BidVolume -= overage
				x.nextBar.BidVolume = overage
			} else if rec.AskVolume >= rec.BidVolume && rec.AskVolume >= overage {
				barRow.AskVolume -= overage
				x.nextBar.AskVolume = overage
			} else {
				if rec.AskVolume > rec.BidVolume {
					barRow.BidVolume -= rec.BidVolume
					x.nextBar.BidVolume = 0
					overage -= rec.BidVolume
					x.nextBar.AskVolume -= overage
					barRow.AskVolume -= overage
				} else {
					barRow.AskVolume -= rec.AskVolume
					x.nextBar.AskVolume = 0
					overage -= rec.AskVolume
					x.nextBar.BidVolume -= overage
					barRow.BidVolume -= overage
				}
			}
			return barRow, nil
		}
	}
	return barRow, nil
}

// Volume bars should never bundle...I think.
func (x *VolumeBarAccumulator) AccumulateBar(r *scid.ScidReader) (Bar, error) {
	var barRow BasicBar

	rec, err := r.NextRecord()
	if err != nil {
		return barRow, err
	}

	if x.nextBar.TotalVolume > 0 {
		barRow = x.nextBar
		x.nextBar = BasicBar{}
	} else {
		barRow = BasicBar{IntradayRecord: *rec}
		barRow.DateTime = rec.DateTimeSC.Time()
		barRow.Open = rec.Close
	}

	for {
		rec, err := r.NextRecord()
		if err != nil {
			return barRow, err
		}

		if rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
			normalizeIndexData(rec)
		}

		if rec.DateTimeSC >= x.scdt_endTime {
			return barRow, io.EOF
		}

		updateBar(&barRow, rec)

		if barRow.TotalVolume < x.barSize {
			continue
		}
		if barRow.TotalVolume == x.barSize {
			x.nextBar = BasicBar{}
			return barRow, nil
		}
		if barRow.TotalVolume > x.barSize {

			overage := barRow.TotalVolume - x.barSize

			x.nextBar = BasicBar{IntradayRecord: *rec}
			x.nextBar.DateTime = rec.DateTimeSC.Time()
			x.nextBar.Open = rec.Close

			barRow.TotalVolume -= overage
			x.nextBar.TotalVolume = overage

			if rec.BidVolume >= rec.AskVolume && rec.BidVolume >= overage {
				barRow.BidVolume -= overage
				x.nextBar.BidVolume = overage
			} else if rec.AskVolume >= rec.BidVolume && rec.AskVolume >= overage {
				barRow.AskVolume -= overage
				x.nextBar.AskVolume = overage
			} else {
				if rec.AskVolume > rec.BidVolume {
					barRow.BidVolume -= rec.BidVolume
					x.nextBar.BidVolume = 0
					overage -= rec.BidVolume
					x.nextBar.AskVolume -= overage
					barRow.AskVolume -= overage
				} else {
					barRow.AskVolume -= rec.AskVolume
					x.nextBar.AskVolume = 0
					overage -= rec.AskVolume
					x.nextBar.BidVolume -= overage
					barRow.BidVolume -= overage
				}
			}
			return barRow, nil
		}
	}
	return barRow, nil
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

func parseBarSize(barSize string) (bartype.BarType, int64) {
	/*
		t, duration := bartype.Parse(barSize)
		var duration int64
		if t != bartype.Time {
			duration, _ = strconv.ParseInt(barSize[0:len(barSize)-1], 10, 64)
		} else {
			d, err := time.ParseDuration(barSize)
			if err != nil {
				return t, duration = int64(d)
			}
		}
		return t, duration
	*/
	return bartype.Parse(barSize)
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
			//log.Tracef("Bundled trade: %s", rec)
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
