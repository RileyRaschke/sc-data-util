package scid

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func Test_scidReaderStdIn(t *testing.T) {

	header := &IntradayHeader{UTCStartIndex: uint32(time.Now().Unix())}
	copy((*header).UniqueHeaderId[:], "SCID")

	tmpfile, err := os.CreateTemp("", "scid-test")
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(header.ToBytes()); err != nil {
		t.Errorf("Failed to write IntradayHeader{} to temp file with error: %v", err)
	}
	tmpfile.Seek(0, 0)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile

	//var r *ScidReader
	r, err := ReaderFromFile(os.Stdin)
	if err != nil {
		t.Errorf("Failed to open empty test file with error: %v", err)
	}
	if err := r.JumpToUnix(time.Now().Unix()); err != io.EOF {
		t.Errorf("Unexpected error from JumpToUnix(): %v", err)
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
