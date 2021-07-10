package util


import (
    "os"
    "bufio"
    log "github.com/sirupsen/logrus"
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

    return bufio.NewWriter( fh ), err
}


