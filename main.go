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
	"github.com/RileyR387/sc-data-util/scid"
	"github.com/pborman/getopt/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Version = "undefined"
)

func main() {
	log.Infof("Symbol: %v\n", *symbol)
	var r *scid.ScidReader
	var err error
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
	if *startUnixTime != 0 {
		r.JumpToUnix(*startUnixTime)
	}
	// Dump Raw Ticks
	if *barSize == "" {
		log.Info("Dumping ticks to stdout")
		csv.DumpRawTicks(os.Stdout, r)
	} else {
		// 15m 1h 2d 4h 32t 3200t
		log.Infof("Dumping %v bars to stdout", *barSize)
		csv.DumpBarCsv(os.Stdout, r, time.Unix(*startUnixTime, 0), time.Unix(*endUnixTime, 0), *barSize)
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
	fmt.Printf(`Usage: %s

Activity log is written to STDERR.

OPTIONS
%s
`, me, u[1])

	os.Exit(1)
}

func ShowVersion() {
	fmt.Printf("%v %v\n", me, Version)
	os.Exit(0)
}
