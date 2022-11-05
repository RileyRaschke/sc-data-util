package util

import (
	"io"
	"time"

	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util/bartype"
	log "github.com/sirupsen/logrus"
)

type BarProfileAccumulator interface {
	AccumulateProfile(*scid.ScidReader) (Bar, BarProfile, error)
}

func NewBarProfileAccumulator(startTime time.Time, endTime time.Time, barSize string, bundleOpt bool, withProfile bool) BarProfileAccumulator {
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

func (x *TimeBarAccumulator) AccumulateProfile(r *scid.ScidReader) (Bar, BarProfile, error) {
	var barRow BasicBar
	var barProfile BarProfile

	if x.nextBar.TotalVolume > 0 {
		barRow = x.nextBar
		barProfile = x.nextProfile
		x.nextBar = BasicBar{}
		x.nextProfile = BarProfile{}
	}
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return barRow, barProfile, err
		}
		if x.bundle && rec.Open == scid.FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			err := bundleTradesWithProfile(r, rec, &barProfile)
			if err != nil {
				log.Warn("Error occured before trade was bundled!")
				return barRow, barProfile, err
			}
		} else if rec.Open != scid.SINGLE_TRADE_WITH_BID_ASK {
			normalizeIndexData(rec)
		}

		if rec.DateTimeSC >= x.scdt_endTime {
			return barRow, barProfile, io.EOF
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

			x.nextProfile = BarProfile{}
			x.nextProfile.AddRecord(rec)
			//x.nextProfile.DateTime = x.scdt_barStart.Time()
			//x.nextProfile.Open = rec.Close
			if barRow.TotalVolume != 0 {
				return barRow, barProfile, nil
			}
		} else {
			updateBarWithProfile(&barRow, &barProfile, rec)
		}
	}
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
