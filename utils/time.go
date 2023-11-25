package utils

import "time"

func DurationUntilNextInterval(now time.Time, interval time.Duration) time.Duration {
	return now.Truncate(interval).Add(interval).Sub(now)
}
