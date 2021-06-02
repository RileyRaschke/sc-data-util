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
    //"github.com/spf13/viper"
    "github.com/RileyR387/sc-data-util/scid"
)

var (
    version = "undefined"
)

func init() {
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    for{
        peekHeader, err := reader.Peek(4)
        if err != nil {
            break;
        }
        if string(peekHeader) == "SCID" {
            headerBytes := make([]byte, 56)
            bytesRead, err := io.ReadFull( reader, headerBytes )
            if bytesRead != 56 || err != nil {
                log.Errorf("Failed to read intraday data: %v", err)
            }
            header := scid.IntradayHeaderFromBytes( headerBytes )
            fmt.Printf("Got header: %v\n", header)
        } else {
            raw_scid_record := make([]byte, 40)
            bytesRead, err := io.ReadFull( reader, raw_scid_record)
            if bytesRead != 40 || err != nil {
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
