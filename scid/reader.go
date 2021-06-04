package scid

import (
    "fmt"
    "os"
    "io"
    "bufio"
    "errors"
    "time"
    log "github.com/sirupsen/logrus"
)

const SCID_HEADER_SIZE_BYTES = int(56)
const SCID_RECORD_SIZE_BYTES = int(40)

type ScidDataReader interface {
    io.ReadWriteSeeker
    NextRecord() (*IntradayRecord)
    ReadSince(time.Time) ([]*IntradayRecord)
    Append([]*IntradayRecord) (error)
}

type ScidReader struct {
    io.Reader
    io.Writer
    io.Seeker
    filePath string
    fileHeader *IntradayHeader
    fileHandle *os.File
}

func ReaderFromFile(file interface{}) (*ScidReader, error){
    var err error
    var fh *os.File
    filePath := ""
    ok := true

    fh, ok = file.(*os.File)
    if !ok {
        filePath = file.(string)
        fh, err = os.Open(filePath)
        if err != nil {
            return nil, err
        }
    } else {
        fInfo, err := fh.Stat()
        if err != nil {
            return nil, err
        }
        filePath = fInfo.Name()
    }
    reader := bufio.NewReader( fh )

    peekHeader, err := reader.Peek(4)
    if err != nil {
        return nil, err
    }
    if string(peekHeader) != "SCID" {
        fmtStr := "Failed to open \"%v\" - \".scid\" header check failed."
        msg := fmt.Sprintf(fmtStr, filePath)
        log.Error(msg)
        return nil, errors.New(msg)
    }

    headerBytes := make([]byte, SCID_HEADER_SIZE_BYTES)
    bytesRead, err := io.ReadFull( reader, headerBytes )
    if err != nil {
        log.Errorf("Failed to open \"%v\" with error: %v", filePath, err)
        return nil, errors.New( fmt.Sprintf("Failed to open \"%v\" with error: %v", filePath, err))
    }
    if bytesRead != SCID_HEADER_SIZE_BYTES {
        fmtStr := "Failed to open \"%v\" - Incomplete file header, read %v bytes, expected %v bytes"
        msg := fmt.Sprintf(fmtStr, filePath, bytesRead, SCID_HEADER_SIZE_BYTES)
        log.Error(msg)
        return nil, errors.New( msg )
    }
    header := IntradayHeaderFromBytes( headerBytes )

    x := &ScidReader{
        Reader: reader,
        Writer: nil,
        Seeker: nil,
        filePath: filePath,
        fileHeader: header,
        fileHandle: fh,
    }
    return x, nil
}

func (sr *ScidReader) AsReader() (io.Reader) {
    return *sr
}

func (sr *ScidReader) NextRecord() (*IntradayRecord, error) {
    raw_scid_record := make([]byte, SCID_RECORD_SIZE_BYTES)
    bytesRead, err := io.ReadFull( sr.Reader, raw_scid_record)
    if err != nil {
        return nil, err
    }
    if bytesRead != SCID_RECORD_SIZE_BYTES || err != nil {
        log.Errorf("Failed to read intraday data with error: %v", err)
    }
    return IntradayRecordFromBytes( raw_scid_record ), nil
}

func (sr *ScidReader) ReadSinceUnixSeconds() ([]*IntradayRecord) {
    return []*IntradayRecord{}
}

func (sr *ScidReader) Append(x []*IntradayRecord) (err error) {
    return nil
}

