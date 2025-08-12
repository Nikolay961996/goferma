package utils

import (
	"fmt"
	"time"
)

func RetryerCon(f func() error, isRetryable func(error) bool) error {
	// 1.2...3.....4x
	tryStep := 1
	for tryStep <= 4 {
		Log.Warn(fmt.Sprintf("retry step: %d", tryStep))
		err := f()
		if err == nil {
			return nil
		}
		if !isRetryable(err) || tryStep == 4 {
			return err
		}
		time.Sleep(time.Duration(tryStep*2-1) * time.Second)
		tryStep++
	}

	return nil
}
