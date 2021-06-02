package scid

import(
    "fmt"
    "time"
)

type SCDateTimeMS int64

const SC_EPOCH_OFFSET = int64(-2209161600)

func (x SCDateTimeMS) UnixTimeMicroSeconds() (int64) {
    return int64(x)+(SC_EPOCH_OFFSET*1000000)
}

func (x SCDateTimeMS) UnixTime() (int64) {
    return (int64(x)/1000000)+SC_EPOCH_OFFSET
}

func (x SCDateTimeMS) Time() (time.Time) {
    secs := x.UnixTime()
    nanoSecs := (int64(x)%1000000)*1000
    return time.Unix(secs, nanoSecs)
}
func (x SCDateTimeMS) String() (string){
    return fmt.Sprintf("%v", x.Time().UTC())
    //return fmt.Sprintf("%v", x.UnixTime() )
    //return fmt.Sprintf("%v", int64(x) )
}
