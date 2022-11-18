package util

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/RileyR387/sc-data-util/scid"
	log "github.com/sirupsen/logrus"
)

var (
	Version = "undefined"
)

func WriteBuffer(outFile interface{}) (*bufio.Writer, error) {
	var err error
	var fh *os.File
	filePath := ""
	ok := true

	fh, ok = outFile.(*os.File)

	if !ok {
		filePath = outFile.(string)
		fh, err = os.Open(filePath)
		if err != nil {
			return nil, err
		}
		log.Infof("Writing to file: %v", filePath)
	} else {
		fInfo, err := fh.Stat()
		if err != nil {
			return nil, err
		}
		filePath = fInfo.Name()
		log.Infof("Writing to %v", filePath)
	}

	return bufio.NewWriter(fh), err
}

func updateBarWithProfile(barRow *BasicBar, barProfile *BarProfile, rec *scid.IntradayRecord) {
	if rec.High > barRow.High {
		barRow.High = rec.High
	}
	if rec.Low < barRow.Low {
		barRow.Low = rec.Low
	}
	barRow.Close = rec.Close
	barRow.NumTrades += rec.NumTrades
	barRow.TotalVolume += rec.TotalVolume
	barRow.BidVolume += rec.BidVolume
	barRow.AskVolume += rec.AskVolume
	barProfile.AddRecord(rec)
}

func updateBarProfile(barProfile *BarProfile, rec *scid.IntradayRecord) {
	barProfile.AddRecord(rec)
}

func bundleTrades(r *scid.ScidReader, bundle *scid.IntradayRecord) error {
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return err
		}
		bundle.TotalVolume += rec.TotalVolume
		bundle.BidVolume += rec.BidVolume
		bundle.AskVolume += rec.AskVolume
		bundle.NumTrades += rec.NumTrades
		if rec.High > bundle.High {
			bundle.High = rec.High
		}
		if rec.Low < bundle.Low {
			bundle.Low = rec.Low
		}
		if rec.Open == scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			// assume the last record is the correct close
			bundle.Close = rec.Close
			//log.Tracef("Bundled trade: %s", rec)
			return nil
		}
	}
	return nil
}

func bundleTradesWithProfile(r *scid.ScidReader, bundle *scid.IntradayRecord, profile *BarProfile) error {
	for {
		rec, err := r.NextRecord()
		if err != nil {
			return err
		}
		profile.AddRecord(rec)
		bundle.TotalVolume += rec.TotalVolume
		bundle.BidVolume += rec.BidVolume
		bundle.AskVolume += rec.AskVolume
		bundle.NumTrades += rec.NumTrades
		if rec.High > bundle.High {
			bundle.High = rec.High
		}
		if rec.Low < bundle.Low {
			bundle.Low = rec.Low
		}
		if rec.Open == scid.LAST_SUB_TRADE_OF_UNBUNDLED_TRADE {
			// assume the last record is the correct close
			bundle.Close = rec.Close
			//log.Tracef("Bundled trade: %s", rec)
			return nil
		}
	}
	return nil
}

func normalizeIndexData(rec *scid.IntradayRecord) {
	// support for index style data
	if rec.High == rec.Low {
		if rec.High < rec.Open {
			log.Debugf("High(%f) is below the Open(%f) at %s", rec.High, rec.Open, rec.DateTimeSC)
			rec.High = rec.Open
		}
		if rec.High < rec.Close {
			log.Debugf("High(%f) is below the Close(%f) at %s", rec.High, rec.Open, rec.DateTimeSC)
			rec.High = rec.Close
		}
		if rec.Low > rec.Open {
			log.Debugf("Low(%f) is above the Open(%f) at %s", rec.Low, rec.Open, rec.DateTimeSC)
			rec.Low = rec.Open
		}
		if rec.Open != 0 && rec.Low < 0.95*rec.Open {
			log.Debugf("Low(%f) is 95%% below Open(%f) at %s", rec.Low, rec.Open, rec.DateTimeSC)
			rec.Low = rec.Open
		}
		if rec.Low > rec.Close {
			log.Debugf("Low(%f) is above the Close(%f) at %s", rec.Low, rec.Close, rec.DateTimeSC)
			rec.Low = rec.Close
		}
		if rec.Close != 0 && rec.Low < 0.95*rec.Close {
			log.Debugf("Low(%f) is 95%% below Close(%f) at %s", rec.Low, rec.Close, rec.DateTimeSC)
			rec.Low = rec.Close
		}
	}
}

func FormatPrice(p float32, tickSizeStr string) string {
	tickSize, formatStr, _ := ParseTickSize(tickSizeStr)
	return fmt.Sprintf(formatStr, RoundToTickSize(p, tickSize))
}

func ParseTickSize(tickSizeStr string) (tickSize float64, formatStr string, err error) {
	tickSize, err = strconv.ParseFloat(strings.TrimSpace(tickSizeStr), 64)
	if err != nil {
		return 0.0, "", err
	}
	formatStr = fmt.Sprintf("%%.%vf", len(tickSizeStr[strings.Index(tickSizeStr, ".")+1:]))
	return
}

func RoundToTickSize(n float32, s float64) float64 {
	size := 1 / s
	return math.Round(float64(n)*size) / size
}
