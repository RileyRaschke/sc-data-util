package dly

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type DayData struct {
	TradingDay  time.Time
	Open        float32
	High        float32
	Low         float32
	Close       float32
	PriorSettle float32
	Volume      int64
}

type DailyDetail struct {
	FilePath  string
	Symbol    string
	MonthCode string
	FirstDay  time.Time
	LastDay   time.Time
}

type DailyData struct {
	io.Reader
	filePath   string
	fileHandle *os.File
	Detail     DailyDetail
}

type DailyDataProvider struct {
	Symbol       string
	DataDir      string
	DailyByMonth map[string]*DailyData
	dataFiles    []string
}

func NewDailyDataProvider(symbol string, dataDir string) *DailyDataProvider {
	dp := &DailyDataProvider{Symbol: strings.ToUpper(symbol), DataDir: dataDir}
	dp.scanMonths()
	return dp
}

func DailyDataFromFile(fPath, symbol string) *DailyData {
	log.Infof("Loading: %v", fPath)

	fName := filepath.Base(fPath)
	monthCode := strings.Replace(fName, symbol, "", -1)
	monthCode = strings.Replace(monthCode, ".dly", "", -1)

	dDetail := DailyDetail{
		fPath,
		symbol,
		monthCode,
		time.Time{},
		time.Time{},
	}

	dd := &DailyData{Detail: dDetail}
	dd.Refresh()
	return dd
}

func (dp *DailyDataProvider) DumpDetailCsv() {
}

func (dp *DailyDataProvider) scanMonths() {
	dp.DailyByMonth = make(map[string]*DailyData)
	err := filepath.Walk(dp.DataDir,
		func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".dly" && !info.IsDir() {
				if strings.HasPrefix(info.Name(), dp.Symbol) {
					monthCode := strings.Replace(info.Name(), dp.Symbol, "", -1)
					monthCode = strings.Replace(monthCode, ".dly", "", -1)
					dp.dataFiles = append(dp.dataFiles, path)
					dd := DailyDataFromFile(path, dp.Symbol)
					dp.DailyByMonth[monthCode] = dd
				}
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
}

func (dd *DailyData) Refresh() error {
	fmt.Printf("Refreshing from: %v\n", dd.Detail.FilePath)
	fh, err := os.Open(dd.Detail.FilePath)
	if err != nil { // we found the file.. shouldn't happen!
		log.Errorf("Failed refresh DailyData from file (%v) with error: %v", dd.Detail.FilePath, err)
		return err
	}

	defer fh.Close()
	reader := bufio.NewReader(fh)
	dd.Reader = reader

	return nil
}

func (dd *DailyData) GetSettle(date time.Time) float32 {
	return 0
}

func (dd *DailyData) GetPriorSettle(date time.Time) float32 {
	return 0
}
