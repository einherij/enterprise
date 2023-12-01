package utils

func Must[T any](v T, err error) T {
	PanicOnError(err)
	return v
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
