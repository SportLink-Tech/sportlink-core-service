package matchannouncement

import (
	"fmt"
	"time"
)

// TimeSlot represents a time range for a match
type TimeSlot struct {
	StartTime time.Time
	EndTime   time.Time
}

func NewTimeSlot(startTime, endTime time.Time) (TimeSlot, error) {
	if endTime.Before(startTime) {
		return TimeSlot{}, fmt.Errorf("end time cannot be before start time")
	}
	return TimeSlot{
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// Duration returns the duration of the slot
func (ts TimeSlot) Duration() time.Duration {
	return ts.EndTime.Sub(ts.StartTime)
}

// Contains checks if a given time is within the slot
func (ts TimeSlot) Contains(t time.Time) bool {
	return (t.Equal(ts.StartTime) || t.After(ts.StartTime)) && t.Before(ts.EndTime)
}
