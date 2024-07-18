package zutil

import (
	"strconv"
	"time"
)

var _secondUnits = []string{"ns", "us", "ms", "s"}

func ManualSecond(du time.Duration) string {
	d := float64(du)

	i := 0
	for d > 1000 && i < len(_secondUnits)-1 {
		d = d / 1000
		i++
	}

	return strconv.FormatFloat(d, 'f', 2, 64) + _secondUnits[i]
}
