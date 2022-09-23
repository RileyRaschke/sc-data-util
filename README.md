
# Sierra Chart Data Utility

Do things with SierraChart data files. Used for automating data analysis with `R` and other things.

## Building

Installing go and executing `go build` should suffice. Version will not be embedded with this approach, but that's mostly for me.

## Command Line tool

Some configuration is required, see usage for `--genconfig` option details.

* Parse and **aggregate** Sierra Chart Intraday Data files (`*.scid`)
* Start and End time filters
* Volume based bar exports
* Tick (trade) based bar exports, bundled and unbundled.
* (coming soon) Support for automated roll/continuous bar exports across expiries
* (coming soon) Expanded CSV columns to support flattening of trading day, prior settlement, and open interest for eased analysis
* (coming soon) cached continuous bars (by size) across expiries, expanded on demand


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
                    Note: T for ticks, V for volume (like Sierra)
                    Note: To export raw ticks (raw times and sales), don't provide this option

     --dailyDetail  Print daily data with added row detail

     --endUnixTime=value
                    End export at unix time [1663892348]

     --startUnixTime=value
                    Export Starting at unix time

 -i, --stdin        Read data from STDIN, Dump to STDOUT. Disables most other
                    options.

 -m, --bundle       Cout a large order filling multiple limits as one bundled trade.

     --slim         Slim/Minimal CSV data

 -s, --symbol=value
                    Symbol to operate on (required, unless `-i`)
     --version      Show version (undefined)
 -x, --genconfig    Write example config to "./sc-data-util.yaml"

```

