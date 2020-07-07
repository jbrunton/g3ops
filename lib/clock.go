package lib

import "time"

// Clock - represents current time
type Clock interface {
	Now() time.Time
}

// SystemClock - a Clock instance using the system time
type SystemClock struct{}

// Now - returns the current time
func (clock *SystemClock) Now() time.Time {
	return time.Now()
}

// NewSystemClock - returns a new concrete clock instance
func NewSystemClock() *SystemClock {
	return &SystemClock{}
}
