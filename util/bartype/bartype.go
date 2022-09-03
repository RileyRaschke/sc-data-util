package bartype

import (
	"strconv"
	"time"
)

type BarType int

const (
	Time BarType = iota
	Tick
	Volume
	Daily
	Weekly
	Monthly
	Quarterly
	Yearly
)

var (
	InvalidBarType = "Invalid Bar Size"
)

func (s BarType) String() string {
	switch s {
	case Time:
		return "Time"
	case Tick:
		return "Tick"
	case Volume:
		return "Volume"
	case Daily:
		return "Daily"
	case Weekly:
		return "Weekly"
	case Monthly:
		return "Monthly"
	case Quarterly:
		return "Quarterly"
	case Yearly:
		return "Yearly"
	default:
		return InvalidBarType
	}
}

func ValidBarType(s BarType) bool {
	if s.String() == InvalidBarType {
		return false
	}
	return true
}

func ParseType(s string) BarType {
	flag := s[len(s)-1:]
	switch flag {
	case "v":
		return Volume
	case "t":
		return Tick
	case "d":
		fallthrough
	case "D":
		return Daily
	case "w":
		fallthrough
	case "W":
		return Weekly
	case "M":
		return Monthly
	case "Q":
		return Quarterly
	case "Y":
		return Yearly
	default:
		return Time
	}
}

func Parse(s string) (BarType, int64) {
	var duration int64
	bt := ParseType(s)
	switch bt {
	case Time:
		d, err := time.ParseDuration(s)
		if err != nil {
			return bt, int64(d)
		}
		break
	case Volume:
		fallthrough
	case Tick:
		duration, _ = strconv.ParseInt(s[0:len(s)-1], 10, 64)
		return bt, duration
	case Daily:
		return bt, -1
	case Weekly:
		return bt, -2
	case Monthly:
		return bt, -3
	case Quarterly:
		return bt, -4
	case Yearly:
		return bt, -5
	}
	return bt, duration
}
