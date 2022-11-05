package util

import (
	"io"
	"time"

	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util/bartype"
	log "github.com/sirupsen/logrus"
)

type BarAccumulator interface {
	AccumulateBar(*scid.ScidReader) (Bar, error)
}

func NewBarAccumulator(startTime time.Time, endTime time.Time, barSize string, bundleOpt bool, withProfile bool) BarAccumulator {
	bt, duration := bartype.Parse(barSize)
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
