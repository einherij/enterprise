package utils

import (
	"strconv"
	"time"
)

func DurationUntilNextInterval(now time.Time, interval time.Duration) time.Duration {
	return now.Truncate(interval).Add(interval).Sub(now)
}

type UnixTime time.Time

func Unix(sec int64) UnixTime {
	return UnixTime(time.Unix(sec, 0))
}

func (ut *UnixTime) Time() time.Time {
	return time.Time(*ut)
}

func (ut *UnixTime) MarshalQuery() (param string, err error) {
	return strconv.Itoa(int(time.Time(*ut).Unix())), nil
}

func (ut *UnixTime) UnmarshalQuery(param string) (err error) {
	i, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return err
	}
	*ut = UnixTime(time.Unix(i, 0))
	return nil
}
