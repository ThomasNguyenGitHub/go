package util

const (
	retryTimes = 5
)

func Retry(rf retryFunc) (v interface{}, retryCount int, err error) {
	retryCount = 0
	for retryCount < retryTimes {
		if v, err = rf(); err == nil {
			return
		}
		retryCount++
	}
	return
}

type retryFunc func() (interface{}, error)
