package zutil

import (
	"time"
)

type Timer struct {
	start time.Time
	pause time.Time
}

func NewTimer() Timer {
	return Timer{
		start: time.Now(),
		pause: time.Now(),
	}
}

func (t Timer) StartAt() time.Time {
	return t.start
}

func (t Timer) PauseAt() time.Time {
	return t.pause
}

func (t Timer) Pause() time.Duration {
	d := time.Now().Sub(t.pause)
	t.pause = time.Now()
	return d
}

func (t Timer) Pauses() string {
	return ManualSecond(t.Pause())
}

func (t Timer) Stop() time.Duration {
	return time.Now().Sub(t.start)
}

func (t Timer) Stops() string {
	return ManualSecond(t.Stop())
}
