package utils

import (
	"github.com/avast/retry-go"
	"go.uber.org/zap"

	"github.com/artnoi43/superwatcher/pkg/logger"
)

func Retry(f func() error) error {
	err := retry.Do(func() error {
		if err := f(); err != nil {
			logger.Debug("error in utils.Retry", zap.Error(err))
			return err
		}
		return nil
	},
		retry.Attempts(30),
		retry.Delay(300),
	)

	return err
}

func RetryWithReturn[T any](f func() (T, error)) (T, error) {
	var t T
	var err error

	err = retry.Do(func() error {
		t, err = f()
		if err != nil {
			logger.Error("error in utils.RetryWithReturn", zap.Error(err))
			return err
		}

		return nil
	},
		retry.Attempts(30),
		retry.Delay(300),
	)

	return t, nil
}
