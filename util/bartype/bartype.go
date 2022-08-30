package bartype

import (
	"strings"
)

type BarType int

const (
	Time BarType = iota
	Tick
	Volume
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
	flag = strings.ToLower(flag)
	switch flag {
	case "v":
		return Volume
	case "t":
		return Tick
	default:
		return Time
	}
}
