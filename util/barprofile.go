package util

import "time"

type BarProfile struct {
	DateTime       time.Time
	ProfileDataMap map[int]*ProfileEntry
}

func (x *BarProfile) String() string {
	return "FIXME"
}
