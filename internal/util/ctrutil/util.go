package ctrutil

func IsEmpty[T any](data []T) bool {
	return len(data) == 0
}

type EqualFunc[T any] func(current T) bool
type EqualFuncMap[K comparable, V any] func(key K, val V) bool

func IsExist[T any](haystack []T, eq EqualFunc[*T]) bool {
	for _, obj := range haystack {
		if eq(&obj) {
			return true
		}
	}
	return false
}

type EqualFunc2[T any, U any] func(haystack T, needle U) bool

func NormalEqualFunc[T comparable](val1 *T, val2 *T) bool {
	return *val1 == *val2
}
func IsExists[T any, U any](haystack []T, needle []U, equalFunc EqualFunc2[*T, *U]) bool {
	for _, obj1 := range needle {
		for _, obj2 := range haystack {
			if !equalFunc(&obj2, &obj1) {
				return false
			}
		}
	}
	return true
}

func MapIsExist[K comparable, V any](haystack map[K]V, equalFunc EqualFuncMap[K, V]) bool {
	for k, v := range haystack {
		if equalFunc(k, v) {
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

type SafeConvertFunc[From any, To any] func(current From) (To, error)
type ConvertFunc[From any, To any] func(current From) To

// SafeConvertSliceType Used to convert []From into []To based on function parameter with error checking
func SafeConvertSliceType[From any, To any](slice []From, convertFunc SafeConvertFunc[*From, To]) ([]To, error) {
	result := make([]To, 0, len(slice))
	for _, val := range slice {
		cur, err := convertFunc(&val)
		if err != nil {
			return nil, err
		}
		result = append(result, cur)
	}
	return result, nil
}

// ConvertSliceType Used to convert []From into []To based on function parameter without error checking, use it when the object conversion is always success
func ConvertSliceType[From any, To any](slice []From, convertFunc ConvertFunc[*From, To]) []To {
	result := make([]To, 0, len(slice))
	for _, val := range slice {
		result = append(result, convertFunc(&val))
	}
	return result
}
