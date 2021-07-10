package scid

import(
    "testing"
    "os"
    "io"
    "time"
    "fmt"
)

func Test_scidReaderStdIn( t *testing.T ){

    header := &IntradayHeader{
        //UniqueHeaderId: [4]byte{},
        HeaderSize: uint32(56),
        RecordSize: uint32(0),
        Version: uint16(0),
        Unused1: uint16(0),
        UTCStartIndex: uint32(time.Now().Unix()),
        Reserve: [36]byte{},
    }
    copy((*header).UniqueHeaderId[:], "SCID")
    //fmt.Println(header.UniqueHeaderId)

    tmpfile, err := os.CreateTemp("","scid-test")
    if err != nil {
        t.Errorf("Failed to create temp file: %v",err)
    }
    defer os.Remove(tmpfile.Name()) // clean up

    headerBytes := header.ToBytes()

    if _, err := tmpfile.Write(headerBytes); err != nil {
        t.Errorf("Failed to write IntradayHeader{} to temp file with error: %v", err)
    }
    //fmt.Printf("Wrore header bytes: %v\n", headerBytes)
    tmpfile.Seek(0, 0)

    oldStdin := os.Stdin
    defer func() { os.Stdin = oldStdin }() // Restore original Stdin

    os.Stdin = tmpfile

    //var r *ScidReader
    r, err := ReaderFromFile( os.Stdin )
    if err != nil {
        t.Errorf("Failed to open empty test file with error: %v", err)
    }
    if err := r.JumpToUnix(time.Now().Unix()); err != io.EOF {
        t.Errorf("Unexpected error from JumpToUnix(): %v",err)
    }
    for {
        _, err := r.NextRecord()
        if err == io.EOF {
            break
        }
        if err != nil {
            t.Errorf("Failed to read empty intraday data file with error: %v", err)
        }
    }
    fmt.Printf("Read an empty scid file!\n")
}


