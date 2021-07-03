package scid

import (
    /*
    "fmt"
    "os"
    "io"
    "bufio"
    "errors"
    "time"
    log "github.com/sirupsen/logrus"
    */
)

func (sr *ScidReader) AggregateBy(bs string) {
}

func (sr *ScidReader) NextBar() (*IntradayRecord, error) {
    return sr.aggregate()
}

func (sr *ScidReader) aggregate() (*IntradayRecord, error) {
    return sr.NextRecord()
}

