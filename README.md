
# Sierra Chart Data Utility

Do things with SierraChart data files. Used for automating data analysis with `R` and other things.

## Building

Installing go and executing `go build` should suffice. Version will not be embedded with this approach, but that's mostly for me.

## Command Line tool

Some configuration is required, see usage for `--genconfig` option details.

* Parse and **aggregate** Sierra Chart Intraday Data files (`*.scid`)
* Start and End time filters
* Tick (trade) based bar exports, bundled and unbundled.
* Volume based bar exports. Don't expect a 100% match with Sierra. May support disabled record splitting in future for better matching.
* (coming soon) Expanded CSV columns to support flattening of trading day, prior settlement, and open interest for eased analysis
* (coming soon) Support for automated roll/continuous bar exports across expiries


## Library
* Buffered Interface
* Seek to specified time with a binary search to minimize trash reads

### Usage
```txt
Try: sc-data-util --genconfig

Usage: sc-data-util [OPTIONS]

Notes:
 - Config (sc-data-util.yaml) can reside in [. ./etc $HOME/.sc-data-util/ $HOME/etc /etc]
 - Data is written to Stdout
 - Activity log is written to Stderr
 - startUnixTime options sets first bar start time

OPTIONS
 -b, --barSize=value
                    Export as bars of size: [10s, 2m, 4h, 3200t, 5000v]
     --dailyDetail  Print daily data with added row detail
     --detailProfile
                    Print daily data with added row detail and bar profile
     --endUnixTime=value
                    End export at unix time [1668991430]
 -i, --stdin        Read data from STDIN, Dump to STDOUT. Disables most other
                    options.
 -m, --bundle
     --slim         Slim/Minimal CSV data
     --startUnixTime=value
                    Export Starting at unix time
 -s, --symbol=value
                    Symbol to operate on (required, unless `-i`)
     --tickSize=value
                    Round to tick size [0.001]
     --version      Show version (v0.0.3-17-g7db6bba-dirty)
 -x, --genconfig    Write example config to "./sc-data-util.yaml"

```

