package main

import (
    "fmt"
    "io"
    "bufio"
    "os"
    "strings"
    //"os/signal"
    //"time"
    //"syscall"
    log "github.com/sirupsen/logrus"
    "github.com/pborman/getopt/v2"
    "github.com/spf13/viper"
    "github.com/RileyR387/sc-data-util/scid"
)

var (
    version = "undefined"
)

func main() {
    fmt.Printf("Symbol: %v\n", *symbol )
    var r *scid.ScidReader
    var err error
    dataFile := viper.GetString("data.dir") + "/" + strings.ToUpper(*symbol) + ".scid"
    if *stdIn {
        r, err = scid.ReaderFromFile( os.Stdin )
    } else {
        r, err = scid.ReaderFromFile( dataFile )
    }
    if err != nil {
        fmt.Printf("Failed to open file '%v' with error: %v", dataFile, err)
        os.Exit(1)
    }
    for {
        rec, err := r.NextRecord()
        if err == io.EOF {
            break
        }
        //rec.TotalVolume += 1
        fmt.Printf("%v\n", rec )
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
