package timing

import (
	"time"
)

func GetCurrent(object string) int {
	var output int
	currentTime := GetCurrentTime()
	switch object {
		case "day":
			output = currentTime.Day()
		case "month":
			output = int(currentTime.Month())
		case "year":
			output = currentTime.Year()
	}
	return output
}

func TimeDurationPassed(start, end time.Time, limit int, format string) bool {
	var diff float64
	if format == "h" {
		diff = end.Sub(start).Hours()
	}
	if diff >= float64(limit) {
		return true
	} else {
		return false
	}
}

func Sleep(units int, format string) {
	time.Sleep(GetTimeDuration(units, format))
}

func GetCurrentTime() time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc)
}

func GetTimeDuration(units int, format string) time.Duration {
	return time.Duration(units * int(GetTimeFormat(format)))
}

func GetTimeFormat(format string) time.Duration {
	var timeFormat time.Duration
	if format == "ms" {
		timeFormat = time.Millisecond
	} else if format == "s" {
		timeFormat = time.Second
	} else if format == "m" {
		timeFormat = time.Minute
	} else if format == "h" {
		timeFormat = time.Hour
	}
	return timeFormat
}

func GetDayTime() string {
	t := time.Now()
	var daytime string
	switch {
	case t.Hour() < 6:
		daytime = "night"
	case t.Hour() < 11:
		daytime = "morning"
	case t.Hour() < 12:
		daytime = "noon"
	case t.Hour() < 17:
		daytime = "afternoon"
	default:
		daytime = "evening"
	}
	return daytime
}