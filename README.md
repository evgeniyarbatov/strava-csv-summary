# strava-csv-summary

Compress CSV file export from Strava to make it useful.

Ths CSV file is create from the Strava export using [strava2csv](https://github.com/evgeniyarbatov/strava2csv) script.

## run

```
go run main.go \
~/Downloads/strava2gpx/gpx-file.csv \
~/Downloads/strava2gpx/gpx-summary.csv
```

## output

```
StartTime,EndTime,Sport,Filename,Duration,Distance,HRMin,HRMedian,HRMax,ElevationMin,ElevationMedian,ElevationMax,CadenceMin,CadenceMedian,CadenceMax,PowerMin,PowerMedian,PowerMax
2020-08-03T22:09:58.398Z,2020-08-03T23:23:42.816Z,Running,4131085890.tcx.gz,4424.418,13563.582086794788,105,142,158,1.372,5.486,14.478,0,79,97,0,282,390
2023-05-15T22:04:50.835Z,2023-05-15T23:04:00.835Z,Running,9742548194.tcx.gz,3550,11212.791281136295,81,148,170,-8.23,-3.2,3.81,0,79,84,0,306,472
2020-12-24T22:23:42.161Z,2020-12-24T23:49:39.196Z,Running,4823154372.tcx.gz,5157.035,15208.744109494473,113,157,167,-118.873,8.077,16.155,0,80,114,0,286,1006
```

