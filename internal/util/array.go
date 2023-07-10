package util

func IsEmpty[T any](data []T) bool {
	return len(data) == 0
}

type EqualFunc[T any] func(current T) bool

func IsExist[T any](haystack []T, eq EqualFunc[T]) bool {
	for _, obj := range haystack {
		if eq(obj) {
			return true
		}
	}
	return false
}

func MapValues[K comparable, V any](maps map[K]V) []V {
	arr := make([]V, 0, len(maps))
	for _, v := range maps {
		arr = append(arr, v)
	}

	return arr
}

func MapKeys[K comparable, V any](maps map[K]V) []K {
	arr := make([]K, 0, len(maps))
	for k, _ := range maps {
		arr = append(arr, k)
	}

	return arr
}
