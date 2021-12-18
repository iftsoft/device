package core

import "time"

// DurationTimer is used to get duration of process
type DurationTimer struct {
	start time.Time
}

// StartTimer saves time of process beginning
func (timer *DurationTimer) StartTimer() {
	timer.start = time.Now()
}

// Nanoseconds returns process duration in nanoseconds
func (timer *DurationTimer) Nanoseconds() int64 {
	if !timer.start.IsZero() {
		delta := time.Now().Sub(timer.start)
		return int64(delta)
	}
	return 0
}

// Microseconds returns process duration in microseconds
func (timer *DurationTimer) Microseconds() int {
	if !timer.start.IsZero() {
		delta := time.Now().Sub(timer.start)
		return int(delta / time.Microsecond)
	}
	return 0
}

// Milliseconds returns process duration in milliseconds
func (timer *DurationTimer) Milliseconds() int {
	if !timer.start.IsZero() {
		delta := time.Now().Sub(timer.start)
		return int(delta / time.Millisecond)
	}
	return 0
}

// Seconds returns process duration in seconds
func (timer *DurationTimer) Seconds() int {
	if !timer.start.IsZero() {
		delta := time.Now().Sub(timer.start)
		return int(delta / time.Second)
	}
	return 0
}
