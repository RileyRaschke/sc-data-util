package util

import (
	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util/bartype"
)

type TimeBarAccumulator struct {
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	scdt_nextBar  scid.SCDateTimeMS
	scdt_duration scid.SCDateTimeMS
	barType       bartype.BarType
	barSize       int64
	bundle        bool
	withProfile   bool
	nextProfile   BarProfile
	nextBar       BasicBar
}

type TickBarAccumulator struct {
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	barSize       uint32
	bundle        bool
	withProfile   bool
	nextProfile   BarProfile
	nextBar       BasicBar
}

type VolumeBarAccumulator struct {
	scdt_barStart scid.SCDateTimeMS
	scdt_endTime  scid.SCDateTimeMS
	barSize       uint32
	withProfile   bool
	nextProfile   BarProfile
	nextBar       BasicBar
}
