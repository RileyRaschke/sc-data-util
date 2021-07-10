
# Sierra Chart Data Utility

Do things with SierraChart data files. Used for automating data analysis with `R` and other things.

## Command Line tool

Some configuration is required, see usage for `--genconfig` option details.

* Parse and **aggregate** Sierra Chart Intraday Data files (`*.scid`)
* Start and End time filters
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
                  Export as bars of size: [10s, 2m, 4h]
     --endUnixTime=value
                  End export at unix time [1625893664]
 -i, --stdin      Read data from STDIN, Dump to STDOUT. Disables most other
                  options.
     --startUnixTime=value
                  Export Starting at unix time
 -s, --symbol=value
                  Symbol to operate on (required, unless `-i`)
     --version    Show version (d875a56-dirty)
 -x, --genconfig  Write example config to "./sc-data-util.yaml"
```

