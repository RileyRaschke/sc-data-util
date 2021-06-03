package scid

import (
    "fmt"
    "os"
    "io"
    "bufio"
    "errors"
    log "github.com/sirupsen/logrus"
)

const SCID_HEADER_SIZE_BYTES = int(56)
const SCID_RECORD_SIZE_BYTES = int(40)

type ScidDataReader interface {
    io.ReadWriteSeeker
    NextRecord() (*IntradayRecord)
    ReadSince() ([]*IntradayRecord)
    Append([]*IntradayRecord) (error)
}

type ScidReader struct {
    io.Reader
    io.Writer
    io.Seeker
    filePath string
    fileHeader *IntradayHeader
}

func NewReader(file string) (*ScidReader, error){
    fh, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    reader := bufio.NewReader( fh )

    peekHeader, err := reader.Peek(4)
    if err != nil {
        return nil, err
    }
    if string(peekHeader) != "SCID" {
        fmtStr := "Failed to open \"%v\" - \".scid\" header check failed."
        msg := fmt.Sprintf(fmtStr, file)
        log.Error(msg)
        return nil, errors.New(msg)
    }

    headerBytes := make([]byte, SCID_HEADER_SIZE_BYTES)
    bytesRead, err := io.ReadFull( reader, headerBytes )
    if err != nil {
        log.Errorf("Failed to open \"%v\" with error: %v", file, err)
        return nil, errors.New( fmt.Sprintf("Failed to open \"%v\" with error: %v", file, err))
    }
    if bytesRead != SCID_HEADER_SIZE_BYTES {
        fmtStr := "Failed to open \"%v\" - Incomplete file header, read %v bytes, expected %v bytes"
        msg := fmt.Sprintf(fmtStr, file, bytesRead, SCID_HEADER_SIZE_BYTES)
        log.Error(msg)
        return nil, errors.New( msg )
    }
    header := IntradayHeaderFromBytes( headerBytes )

    x := &ScidReader{
        Reader: reader,
        Writer: nil,
        Seeker: nil,
        filePath: file,
        fileHeader: header,
    }
    return x, nil
}

func (sr *ScidReader) AsReader() (io.Reader) {
    return *sr
}

func (sr *ScidReader) Append(x []*IntradayRecord) (err error) {
    return nil
}

func (sr *ScidReader) NextRecord() (*IntradayRecord) {
    raw_scid_record := make([]byte, SCID_RECORD_SIZE_BYTES)
    bytesRead, err := io.ReadFull( sr.Reader, raw_scid_record)
    if bytesRead != SCID_RECORD_SIZE_BYTES || err != nil {
        log.Errorf("Failed to read intraday data with error: %v", err)
    }
    return IntradayRecordFromBytes( raw_scid_record )
}

