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

func init() {
}

func main() {
    fmt.Printf("Symbol: %v\n", *symbol )

    dataFile := viper.GetString("data.dir") + "/" + strings.ToUpper(*symbol) + ".scid"

    r, err := scid.ReaderFromFile( dataFile )
    if err != nil {
        fmt.Printf("Failed to open file '%v' with error: %v", dataFile, err)
        os.Exit(1)
    }
    for {
        rec, err := r.NextRecord()
        if err == io.EOF {
            break
        }
        fmt.Printf("%v\n", rec )
    }
}

func fromStdIn() {
    reader := bufio.NewReader(os.Stdin)
    for{
        peekHeader, err := reader.Peek(4)
        if err != nil {
            break;
        }
        if string(peekHeader) == "SCID" {
            headerBytes := make([]byte, scid.SCID_HEADER_SIZE_BYTES)
            bytesRead, err := io.ReadFull( reader, headerBytes )
            if bytesRead != scid.SCID_HEADER_SIZE_BYTES || err != nil {
                log.Errorf("Failed to read intraday data: %v", err)
            }
            header := scid.IntradayHeaderFromBytes( headerBytes )
            fmt.Printf("Got header: %v\n", header)
        } else {
            raw_scid_record := make([]byte, scid.SCID_RECORD_SIZE_BYTES)
            bytesRead, err := io.ReadFull( reader, raw_scid_record)
            if bytesRead != scid.SCID_RECORD_SIZE_BYTES || err != nil {
                log.Errorf("Failed to read intraday data: %v", err)
            }
            idRec := scid.IntradayRecordFromBytes( raw_scid_record )
            fmt.Printf("%v\n", idRec)
        }
    }
    os.Exit(0)
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
