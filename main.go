package main

import (
	"fmt"
	//"io"
	//"bufio"
	"os"
	"strings"

	//"os/signal"
	"time"
	//"syscall"
	"github.com/RileyR387/sc-data-util/csv"
	"github.com/RileyR387/sc-data-util/dly"
	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Version = util.Version
)

func main() {
	log.Infof("%v version: %v", me, Version)
	log.Infof("Symbol: %v\n", *symbol)
	var r *scid.ScidReader
	var err error

	if *dailyDetail {
		// TODO: Sanitize input?
		dd := dly.NewDailyDataProvider(*symbol, viper.GetString("data.dir"))
		dd.WriteDailyDetailCsv(os.Stdout)
		return
	}

	if *stdIn {
		r, err = scid.ReaderFromFile(os.Stdin)
		if err != nil {
			fmt.Printf("Failed to open os.Stdin with error: %v", err)
			os.Exit(1)
		}
	} else {
		// TODO: test for specified file first? Then ext?
		dataFile := viper.GetString("data.dir") + "/" + strings.ToUpper(*symbol) + ".scid"
		r, err = scid.ReaderFromFile(dataFile)
		if err != nil {
			fmt.Printf("Failed to open file '%v' with error: %v", dataFile, err)
			os.Exit(1)
		}
	}
	// Raw Ticks
	if *barSize == "" {
		log.Info("Writing ticks to stdout")
		csv.WriteRawTicks(os.Stdout, r, time.Unix(*startUnixTime, 0), time.Unix(*endUnixTime, 0), 1)
	} else {
		// 15m 1h 2d 4h 32t 3200t
		// TODO: Support for days/weeks...
		log.Infof("Writing %v bars to stdout", *barSize)
		if *slim {
			csv.WriteBarCsv(os.Stdout, r, time.Unix(*startUnixTime, 0), time.Unix(*endUnixTime, 0), *barSize, *bundle)
		} else if *detailProfile {
			csv.WriteBarDetailWithProfileCsv(os.Stdout, r, time.Unix(*startUnixTime, 0), time.Unix(*endUnixTime, 0), *barSize, *bundle)
		} else {
			csv.WriteBarDetailCsv(os.Stdout, r, time.Unix(*startUnixTime, 0), time.Unix(*endUnixTime, 0), *barSize, *bundle)
		}
	}
}

func usage(msg ...string) {
	if len(msg) > 0 {
		fmt.Fprintf(os.Stderr, "%s\n", msg[0])
	}
	// strip off the first line of generated usage
	b := &strings.Builder{}
	getopt.PrintUsage(b)
	u := strings.SplitAfterN(b.String(), "\n", 2)
	fmt.Printf(`Usage: %s [OPTIONS]

Notes:
 - Config (%v) can reside in %v
 - Data is written to Stdout
 - Activity log is written to Stderr
 - startUnixTime options sets first bar start time

OPTIONS
%s
`, me, yamlFile, configSearchPaths, u[1])

	os.Exit(1)
}

func ShowVersion() {
	fmt.Printf("%v %v\n", me, Version)
	os.Exit(0)
}
