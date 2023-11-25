package utils

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_RetryWithEarlyExitBecauseOfNilErrorReturn(t *testing.T) {
	counter := 0
	callbackFunc := func() error {
		counter++
		if counter == 3 {
			return nil
		}
		return fmt.Errorf("some error")
	}
	Retry(callbackFunc, 5)
	assert.Equal(t, counter, 3)
}

func Test_RetryWithFullTimes(t *testing.T) {
	counter := 0
	callbackFunc := func() error {
		counter++
		return fmt.Errorf("some error")
	}
	Retry(callbackFunc, 10)
	assert.Equal(t, counter, 10)
}

func Test_RetryOverTime(t *testing.T) {
	a := assert.New(t)
	var retries = 4
	getBanner := func() (string, error) {
		retries--
		if retries <= 0 {
			return "banner content", nil
		} else {
			return "", errors.New("banner doesn't exist")
		}
	}
	var banner string
	err := RetryOverTime(func() error {
		var err error
		banner, err = getBanner()
		return err
	}, time.Second, time.Millisecond*50)
	a.NoError(err)
	a.NotEmpty(banner)
}
