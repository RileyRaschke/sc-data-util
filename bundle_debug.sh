#!/bin/bash

make install >&2 || exit 1

symbol=${1:-"f.us.epu22"}
period=${2:-"1000t"}
daysBack=${3:-"22"}

#startTime=$(date +%s -d last-month)
#startTime=$(($(date --date="$(date +%x -d last-month) 17:00:00" +"%s")+(86400*7)))
#endTime=$(($(date --date="$(date +%x) 17:00:00" +"%s")-(86400*$(($daysBack-3)))))

## Sweep down
startTime=$(($(date --date="08/30/2022 09:04:00" +"%s")))
endTime=$(($(date --date="08/30/2022 09:05:30" +"%s")))

## Sweep up
startTime=$(($(date --date="08/31/2022 08:44:00" +"%s")))
endTime=$(($(date --date="08/31/2022 08:45:30" +"%s")))

echo "sc-data-util -b $period --startUnixTime $startTime --endUnixTime $endTime -s $symbol" >&2
sc-data-util -b $period --startUnixTime $startTime --endUnixTime $endTime -s $symbol | column -s, -t

#go build && time ./sc-data-util -s f.us.mnqm21 --startUnixTime $(date --date='05/18/2021 17:00:00' +"%s") --barSize 4h > out.csv && head out.csv && tail out.csv

#go build && time ./sc-data-util -s f.us.mnqm21 --startUnixTime $(date --date='05/18/2021 17:00:00' +"%s") --barSize 15s > out.csv && head out.csv && tail out.csv


