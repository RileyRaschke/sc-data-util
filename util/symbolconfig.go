package util

import (
	"github.com/RileyR387/sc-data-util/util/symbolclass"
)

type SymbolConfig struct {
	symbol               string
	symbolClass          symbolclass.SymbolClass
	sessionStartTime     string
	sessionEndTime       string
	cashStartTime        string
	cashEndTime          string
	newBarAtSessionStart bool
}
