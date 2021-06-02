package main

import (
    "os"
    log "github.com/sirupsen/logrus"
    "runtime"
    "path"
    "strconv"
)

func initLogger(logLevel string, logFile string) {
    var formatter log.Formatter

    level, _ := log.ParseLevel( logLevel )

    log.SetLevel( level )

    if logFile != "" {
        //f, err := os.Create(logFile)
        f, err := os.OpenFile(logFile, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0640)

        if err != nil {
            panic( "Failed to create log file: " + logFile )
        }
        log.SetOutput( f )
        formatter = &log.JSONFormatter{
            TimestampFormat: "2006-01-02 15:04:05.000",
        }
        if log.GetLevel() == log.TraceLevel {
            log.SetReportCaller(true)
            formatter.(*log.JSONFormatter).CallerPrettyfier = LogPrettyTrace
        }
    } else {
        formatter = &log.TextFormatter{
            FullTimestamp: true,
            TimestampFormat: "2006-01-02 15:04:05.000",
            PadLevelText: true,
            QuoteEmptyFields: true,
        }
        if log.GetLevel() == log.TraceLevel {
            log.SetReportCaller(true)
            formatter.(*log.TextFormatter).CallerPrettyfier = LogPrettyTrace
        }
    }
    log.SetFormatter( formatter )
}

func LogPrettyTrace(f *runtime.Frame) (function string, file string){
    //function = path.Base(f.Function)
    function = ""
    file = " " + path.Base(f.File) + ":" + strconv.Itoa(f.Line)
    return function, file
}
