package scid

/**
* Spec: https://www.sierrachart.com/index.php?page=doc/IntradayDataFileFormat.html
 */

import (
	"bytes"
	"encoding/binary"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const SCID_RECORD_SIZE_BYTES = int(40)

// Not sure about these constants yet...
const SINGLE_TRADE_WITH_BID_ASK = float32(0.0)
const FIRST_SUB_TRADE_OF_UNBUNDLED_TRADE = float32(-1.99900095e+37)
const LAST_SUB_TRADE_OF_UNBUNDLED_TRADE = float32(-1.99900197e+37)

/*
* 40 bytes total (320 bits)
 */
type IntradayRecord struct {
	DateTimeSC SCDateTimeMS // 8

	Open  float32 // 4
	High  float32 // 4
	Low   float32 // 4
	Close float32 // 4

	NumTrades   uint32 // 4
	TotalVolume uint32 // 4
	BidVolume   uint32 //4
	AskVolume   uint32 //4
}

func IntradayRecordFromBytes(b []byte) (x *IntradayRecord) {
	x = &IntradayRecord{}
	if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, x); err != nil {
		log.Error("binary.Read failed: ", err)
	}
	return x
}

func (x *IntradayRecord) String() string {
	return fmt.Sprintf("{"+
		"\"t\":\"%v\", "+
		"\"o\":\"%v\", "+
		"\"h\":\"%v\", "+
		"\"l\":\"%v\", "+
		"\"c\":\"%v\", "+
		"\"nt\":\"%v\", "+
		"\"tv\":\"%v\", "+
		"\"bv\":\"%v\", "+
		"\"av\":\"%v\""+
		"}", x.DateTimeSC,
		x.Open,
		x.High,
		x.Low,
		x.Close,
		x.NumTrades,
		x.TotalVolume,
		x.BidVolume,
		x.AskVolume,
	)
}
