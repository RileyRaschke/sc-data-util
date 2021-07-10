
# Sierra Chart Data Utility

## Library

* Buffered Interface
* Partial support for Stdin

## Command Line tool  

* Parse and aggregate Sierra Chart Intraday Data (`*.scid`)
* Seek to specified time with a binary search to minimize reads

### Usage
```
./sc-data-util

No input or symbol provided

Try: sc-data-util --genconfig

Usage: sc-data-util

Activity log is written to STDERR.

OPTIONS
 -b, --barSize=value
                  Export as bars
     --endUnixTime=value
                  End export at unix time [1625889403]
 -i, --stdin      Read data from STDIN, Dump to STDOUT. Disables most other
                  options.
     --json       1m
     --startUnixTime=value
                  Export Starting at unix time
 -s, --symbol=value
                  Symbol to operate on (required, unless `-i`)
     --version    Show version (fb17ee4)
 -x, --genconfig  Write example config to "./sc-data-util.yaml"

```

