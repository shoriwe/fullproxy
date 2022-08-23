package common

type Hashable interface {
	string | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func CopyMap[T Hashable, Q any](m map[T]Q) map[T]Q {
	result := make(map[T]Q, len(m))
	for key, value := range m {
		result[key] = value
	}
	return result
}
