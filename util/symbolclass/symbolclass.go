package symbolclass

import "strings"

type SymbolClass int

const (
	Equity SymbolClass = iota
	Energy
	Rate
	Treasury
	FX
	Metal
	Soft
)

var (
	InvalidSymbolClass = "Invalid Symbol Class"
	AllMonthCodes      = []string{"F", "G", "H", "J", "K", "M", "N", "Q", "U", "V", "X", "Z"}
	EquityMonthCodes   = []string{"H", "M", "U", "Z"}
)

func (s SymbolClass) String() string {
	switch s {
	case Equity:
		return "Equity"
	case Energy:
		return "Energy"
	case Rate:
		return "Rate"
	case Treasury:
		return "Treasury"
	case FX:
		return "FX"
	case Metal:
		return "Metal"
	case Soft:
		return "Soft"
	default:
		return InvalidSymbolClass
	}
}

func ValidSymbolClass(s SymbolClass) bool {
	if s.String() == InvalidSymbolClass {
		return false
	}
	return true
}

func ParseType(s string) SymbolClass {
	flag := strings.ToLower(s)
	switch flag {
	case "equity":
		return Equity
	case "energy":
		return Energy
	case "rate":
		return Rate
	case "treasury":
		return Treasury
	case "fx":
		return FX
	case "metal":
		return Metal
	case "soft":
		return Soft
	default:
		return Rate
	}
}

func (c SymbolClass) MonthCodes() []string {
	switch c {
	case Equity:
		return EquityMonthCodes
	case Energy:
		return AllMonthCodes
	case Rate:
		return AllMonthCodes
	case Treasury:
		return EquityMonthCodes
	case FX:
		return EquityMonthCodes
	case Metal:
		return AllMonthCodes // FIXME
	case Soft:
		return AllMonthCodes // FIXME
	default:
		return AllMonthCodes
	}
}
