package utils

import (
	"errors"
	"time"
)

func Retry(f func() error, times int) {
	for i := 0; i < times; i++ {
		if err := f(); err == nil {
			break
		}
	}
}

func RetryOverTime(f func() error, timeout time.Duration, tick time.Duration) error {
	err := f()
	if err == nil {
		return nil
	}
	success := make(chan bool, 1)
	timer := time.NewTimer(timeout)
	ticker := time.NewTicker(tick)
	for tick := ticker.C; ; {
		select {
		case <-timer.C:
			return errors.New("retry timeout exceeded")
		case <-tick:
			tick = nil
			go func() {
				err := f()
				success <- err == nil
			}()
		case succ := <-success:
			if succ {
				return nil
			}
			tick = ticker.C
		}
	}
}
