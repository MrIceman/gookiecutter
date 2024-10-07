package collectionoperation

func GroupBy[T any, U comparable](arr []T, f func(T) U) map[U][]T {
	result := make(map[U][]T)
	for _, v := range arr {
		key := f(v)
		result[key] = append(result[key], v)
	}
	return result
}
