package util

import (
	"time"
)

var TimeFormat string = "02-01-2006 15:04"
var DateFormat string = "02-01-2006"
var MysqlDateFormat string = "2006-01-02"

func Now() int64 {
	return time.Now().Unix()
}

func Int64ToTimeStr(epoch int64) string {
	rv := time.Unix(epoch, 0)
	return rv.UTC().Format(TimeFormat)
}

func TimeStrToInt64(timeStr string) (int64, error) {
	rv, err := time.Parse(TimeFormat, timeStr)
	if err != nil {
		return 0, err
	}

	return rv.Unix(), nil
}

func IsValidTimeStr(timeStr string) bool {
	epoch, err := TimeStrToInt64(timeStr)
	if err != nil {
		return false
	}

	str := Int64ToTimeStr(epoch)

	return timeStr == str
}

func Int64ToDateStr(epoch int64) string {
	rv := time.Unix(epoch, 0)
	return rv.UTC().Format(DateFormat)
}

func DateStrToInt64(timeStr string) (int64, error) {
	rv, err := time.Parse(DateFormat, timeStr)
	if err != nil {
		return 0, err
	}

	return rv.Unix(), nil
}

func IsValidDateStr(dateStr string) bool {
	epoch, err := DateStrToInt64(dateStr)
	if err != nil {
		return false
	}

	str := Int64ToDateStr(epoch)

	return dateStr == str
}

type clockTimeTicker struct {
	C chan time.Time
}

func ClockTimeTicker(interval time.Duration, offset time.Duration) *clockTimeTicker {
	s := &clockTimeTicker{
		C: make(chan time.Time),
	}

	go func() {
		now := time.Now()

		// Figure out when the first tick should happen
		firstTick := now.Truncate(interval).Add(interval).Add(offset)

		// Block until the first tick
		<-time.After(firstTick.Sub(now))

		t := time.NewTicker(interval)

		// Send initial tick
		s.C <- firstTick

		for {
			// Forward ticks from the native time.Ticker to the ClockTimeTicker channel
			s.C <- <-t.C
		}
	}()

	return s
}
