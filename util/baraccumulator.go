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
		x.bundle = bundleOpt
		x.withProfile = withProfile
		log.Infof("BAR_SIZE %s, DURATION %d, BSTART %s, BNEXTSTART %s", barSize, duration, x.scdt_barStart, x.scdt_nextBar)
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
		}
		/*
			 else if sym.IsIndexData() && rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
				normalizeIndexData(rec)
			}
		*/

		if rec.DateTimeSC >= x.scdt_endTime {
			return barRow, io.EOF
		}

		var i int = 0
		if rec.DateTimeSC >= x.scdt_nextBar {
			x.scdt_barStart = x.scdt_nextBar
			// assure the next tick is within the next bar's duration
			for {
				//log.Infof("recTime %s, nextBar %s endTime %s %d", rec.DateTimeSC, x.scdt_nextBar, x.scdt_endTime, i)
				if x.scdt_nextBar > rec.DateTimeSC {
					//log.Info("Broke out")
					break
				} else {
					x.scdt_barStart = x.scdt_nextBar
					x.scdt_nextBar += x.scdt_duration
				}
				i++
			}
			x.nextBar = NewBasicBar(rec, x.scdt_barStart)

			if barRow.TotalVolume != 0 {
				return barRow, nil
			}
		} else {
			if barRow.IntradayRecord.DateTimeSC < x.scdt_barStart {
				barRow = NewBasicBar(rec, x.scdt_barStart)
			} else {
				barRow.AddRecord(rec)
			}
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
		barRow = NewBasicBar(rec, rec.DateTimeSC)
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
		}
		/*
			 else if sym.IsIndexData() && rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
				normalizeIndexData(rec)
			}
		*/

		if rec.DateTimeSC >= x.scdt_endTime {
			return barRow, io.EOF
		}

		barRow.AddRecord(rec)

		if barRow.NumTrades < x.barSize {
			continue
		}
		if barRow.NumTrades == x.barSize {
			x.nextBar = BasicBar{}
			return barRow, nil
		}
		if barRow.NumTrades > x.barSize {
			overage := barRow.NumTrades - x.barSize

			x.nextBar = NewBasicBar(rec, rec.DateTimeSC)

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
		barRow = NewBasicBar(rec, rec.DateTimeSC)
	}

	for {
		rec, err := r.NextRecord()
		if err != nil {
			return barRow, err
		}
		/*
			if sym.IsIndexData() && rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
				normalizeIndexData(rec)
			}
		*/
		if rec.BidVolume > 0 && rec.AskVolume > 0 {
			log.Debugf("Both bid and ask volume on tick: %v", rec)
		}

		if rec.NumTrades > 1 {
			log.Debugf("More than one trade on tick: %v", rec)
		}

		if rec.DateTimeSC >= x.scdt_endTime {
			return barRow, io.EOF
		}

		barRow.AddRecord(rec)

		if barRow.TotalVolume < x.barSize {
			continue
		}
		if barRow.TotalVolume == x.barSize {
			x.nextBar = BasicBar{}
			return barRow, nil
		}
		if barRow.TotalVolume > x.barSize {

			overage := barRow.TotalVolume - x.barSize

			x.nextBar = NewBasicBar(rec, rec.DateTimeSC)

			barRow.TotalVolume -= overage
			x.nextBar.TotalVolume = overage

			if overage > rec.TotalVolume/2 {
				barRow.NumTrades -= 1
				x.nextBar.NumTrades += 1
			} else if overage == rec.TotalVolume/2 {
				barRow.NumTrades -= 1
				//x.nextBar.NumTrades += 1
			} else {
				x.nextBar.NumTrades -= 1
			}

			if rec.BidVolume > rec.AskVolume {
				log.Debugf("Bid Junk with overage %v ran on tick: %v", overage, rec)
				barRow.BidVolume -= overage
				x.nextBar.BidVolume = overage
			} else {
				log.Debugf("Ask Junk with overage %v ran on tick: %v", overage, rec)
				barRow.AskVolume -= overage
				x.nextBar.AskVolume = overage
			}
			return barRow, nil
		}
	}
	return barRow, nil
}
