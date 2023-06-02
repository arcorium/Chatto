package util

func IsEmpty[T any](data []T) bool {
	return len(data) == 0
}
