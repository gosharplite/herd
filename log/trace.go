package log

import (
	"strings"
	"time"
)

func Trace(s string) (string, time.Time) {
	return s, time.Now()
}

func Un(s string, startTime time.Time) {

	endTime := time.Now()

	ms := float64(endTime.Sub(startTime).Nanoseconds()/100000) / 10.0

	l := 33 - len(s)
	if l < 0 {
		l = 0
	}

	s = s + strings.Repeat(" ", l)

	Info("*** "+s+",\t%v ms", ms)
}
