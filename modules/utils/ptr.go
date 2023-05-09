package utils

import "time"

// TimeFromPtr returns time.Time from *time.Time. If t is nil, returns time.Time{}.
func TimeFromPtr(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return *t
}
