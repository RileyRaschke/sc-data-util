package util

import (
	"time"

	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util/bartype"
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
