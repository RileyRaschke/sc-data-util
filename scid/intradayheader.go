package scid

/**
* Spec: https://www.sierrachart.com/index.php?page=doc/IntradayDataFileFormat.html
 */

import (
	"bytes"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
)

const SCID_HEADER_SIZE_BYTES = int(56)

/**
* 56 bytes == 448 bits
* 1 byte = 8 bits
 */
type IntradayHeader struct {
	UniqueHeaderId [4]byte  // 4
	HeaderSize     uint32   // 4
	RecordSize     uint32   // 4
	Version        uint16   // 2
	Unused1        uint16   // 2
	UTCStartIndex  uint32   // 4
	Reserve        [36]byte //36
}

func IntradayHeaderFromBytes(b []byte) (x *IntradayHeader) {
	x = &IntradayHeader{}
	if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, x); err != nil {
		log.Error("binary.Read failed: ", err)
	}
	return x
}

func (h *IntradayHeader) ToBytes() []byte {
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.LittleEndian, h)
	return bin_buf.Bytes()
}
