package domain

import (
	"fmt"
	"time"
)

type TimeRange struct {
	start time.Time
	end   time.Time
}

func NewTimeRange(start, end time.Time) (TimeRange, error) {
	if !end.After(start) {
		return TimeRange{}, fmt.Errorf(
			"invalid time range: start=%v, end=%v", start, end,
		)
	}
	return TimeRange{
		start: start,
		end:   end,
	}, nil
}

func (t TimeRange) Start() time.Time { return t.start }
func (t TimeRange) End() time.Time   { return t.end }
