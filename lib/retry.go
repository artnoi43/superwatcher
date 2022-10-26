package lib

import "github.com/avast/retry-go"

func Retry(f func() error) error {
	err := retry.Do(func() error {
		return f()
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
			return err
		}

		return nil
	},
		retry.Attempts(30),
		retry.Delay(300),
	)

	return t, nil
}
