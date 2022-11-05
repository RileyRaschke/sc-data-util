package util

import (
	"time"

	"github.com/RileyR387/sc-data-util/scid"
)

type BarProfile struct {
	DateTime       time.Time
	ProfileDataMap map[int]*ProfileEntry
}

func (x *BarProfile) String() string {
	return "FIXME"
}

func (x *BarProfile) AddRecord(rec *scid.IntradayRecord) {

}
