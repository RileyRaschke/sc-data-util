package csv

import (
	"fmt"
	"io"
	"time"

	"github.com/RileyR387/sc-data-util/scid"
	"github.com/RileyR387/sc-data-util/util"
	log "github.com/sirupsen/logrus"
)

const CSV_HEADER_RAW = string("Date,Time,Open,Ask,Bid,Last,Volume,NumTrades,BidVolume,AskVolume")
const CSV_HEADER = string("Date,Time,Open,High,Low,Last,Volume,NumTrades,BidVolume,AskVolume")
const CSV_HEADER_DETAIL = string("Date,Time,Open,High,Low,Last,Volume,NumTrades,BidVolume,AskVolume,PriorLast,PriorSettle,TradingDate,TradingDateTime")
const CSV_HEADER_PROFILE = string("Date,Time,Open,High,Low,Last,Volume,NumTrades,BidVolume,AskVolume,PriorLast,PriorSettle,TradingDate,TradingDateTime,BarProfile")

type WriterOptions struct {
	Writer      interface{}
	Reader      *scid.ScidReader
	StartTime   time.Time
	EndTime     time.Time
	TickSizeStr string
	BarSize     string
	Bundle      bool
}

type CsvBarRow struct {
	util.BasicBar
	PriorSettle float32
	PriorLast   float32
	TradingDate time.Time
	TickSize    float64
	FloatFmt    string
}

type CsvProfileBarRow struct {
	CsvBarRow
	util.BarProfile
}

func (x CsvBarRow) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		x.DateTime.Format("2006/1/02"),
		x.DateTime.Format("15:04:05.000"),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.Open, x.TickSize)),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.High, x.TickSize)),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.Low, x.TickSize)),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.Close, x.TickSize)),
		x.TotalVolume,
		x.NumTrades,
		x.BidVolume,
		x.AskVolume,
	)
}
func (x CsvBarRow) DetailString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		x.DateTime.Format("2006/01/02"),
		x.DateTime.Format("15:04:05.000"),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.Open, x.TickSize)),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.High, x.TickSize)),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.Low, x.TickSize)),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.Close, x.TickSize)),
		x.TotalVolume,
		x.NumTrades,
		x.BidVolume,
		x.AskVolume,
		x.PriorLast,
		x.PriorSettle,
		x.TradingDate.Format("2006/01/02"),
		x.TradingDate.Format("15:04:05.000"),
	)
}

func (x CsvProfileBarRow) DetailProfileString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v",
		x.DateTime.Format("2006/01/02"),
		x.DateTime.Format("15:04:05.000"),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.Open, x.TickSize)),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.High, x.TickSize)),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.Low, x.TickSize)),
		fmt.Sprintf(x.FloatFmt, util.RoundToTickSize(x.Close, x.TickSize)),
		x.TotalVolume,
		x.NumTrades,
		x.BidVolume,
		x.AskVolume,
		x.PriorLast,
		x.PriorSettle,
		x.TradingDate.Format("2006/01/02"),
		x.TradingDate.Format("15:04:05.000"),
		x.BarProfile,
	)
}

func WriteBarCsv(opts WriterOptions) error {
	tickSize, formatStr, _ := util.ParseTickSize(opts.TickSizeStr)
	opts.Reader.JumpTo(opts.StartTime)
	w, err := util.WriteBuffer(opts.Writer)
	if err != nil {
		log.Errorf("Failed to open \"%v\" for writing with error: %v", opts.Writer, err)
	}
	w.WriteString(CSV_HEADER + "\n")
	ba := util.NewBarAccumulator(opts.StartTime, opts.EndTime, opts.BarSize, opts.Bundle, false)
	for {
		bar, err := ba.AccumulateBar(opts.Reader)
		barRow := CsvBarRow{BasicBar: bar.(util.BasicBar), TickSize: tickSize, FloatFmt: formatStr}
		if barRow.TotalVolume != 0 {
			w.WriteString(barRow.String() + "\n")
		}
		if err != nil {
			break
		}
	}
	w.Flush()
	return nil
}

func WriteBarDetailCsv(opts WriterOptions) error {
	tickSize, formatStr, _ := util.ParseTickSize(opts.TickSizeStr)
	opts.Reader.JumpTo(opts.StartTime)
	w, err := util.WriteBuffer(opts.Writer)
	log.Info("Writing detail csv")
	if err != nil {
		log.Errorf("Failed to open \"%v\" for writing with error: %v", opts.Writer, err)
	}
	w.WriteString(CSV_HEADER_DETAIL + "\n")
	ba := util.NewBarAccumulator(opts.StartTime, opts.EndTime, opts.BarSize, opts.Bundle, false)
	for {
		bar, err := ba.AccumulateBar(opts.Reader)
		barRow := CsvBarRow{BasicBar: bar.(util.BasicBar), TickSize: tickSize, FloatFmt: formatStr}
		if barRow.TotalVolume != 0 {
			barRow.TradingDate = barRow.DateTime.Add(time.Hour * 7)
			w.WriteString(barRow.DetailString() + "\n")
		}
		if err != nil {
			break
		}
	}
	w.Flush()
	return nil
}

func WriteBarDetailWithProfileCsv(opts WriterOptions) error {
	tickSize, formatStr, _ := util.ParseTickSize(opts.TickSizeStr)
	opts.Reader.JumpTo(opts.StartTime)
	w, err := util.WriteBuffer(opts.Writer)
	log.Info("Writing detail csv with profile")
	if err != nil {
		log.Errorf("Failed to open \"%v\" for writing with error: %v", opts.Writer, err)
		return err
	}
	w.WriteString(CSV_HEADER_PROFILE + "\n")
	ba := util.NewBarProfileAccumulator(opts.StartTime, opts.EndTime, opts.BarSize, opts.Bundle, true)
	for {
		bar, pro, err := ba.AccumulateProfile(opts.Reader)
		br := CsvBarRow{BasicBar: bar.(util.BasicBar), TickSize: tickSize, FloatFmt: formatStr}
		barRow := CsvProfileBarRow{CsvBarRow: br, BarProfile: pro}
		if barRow.TotalVolume != 0 {
			barRow.TradingDate = barRow.DateTime.Add(time.Hour * 7)
			w.WriteString(barRow.DetailProfileString() + "\n")
		}
		if err != nil {
			break
		}
	}
	w.Flush()
	return nil
}

func WriteRawTicks(opts WriterOptions) error {
	opts.Reader.JumpTo(opts.StartTime)
	w, err := util.WriteBuffer(opts.Writer)
	scdt_endTime := scid.NewSCDateTimeMS(opts.EndTime)
	if err != nil {
		log.Errorf("Failed to open \"%v\" for writing with error: %v", opts.Writer, err)
		return err
	}
	w.WriteString(CSV_HEADER_RAW + "\n")
	for {
		rec, err := opts.Reader.NextRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Infof("Error returned by `r.NextRecord()`: %v", err)
		}
		if rec.DateTimeSC >= scdt_endTime {
			break
		}
		//barRow := CsvBarRow{BasicBar: util.BasicBar{IntradayRecord: *rec}}
		//w.WriteString(fmt.Sprintf("%v\n", barRow.TickString()))
		w.WriteString(fmt.Sprintf("%s\n", rec))
	}
	w.Flush()
	return nil
}
