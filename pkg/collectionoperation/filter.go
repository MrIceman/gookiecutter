package collectionoperation

func Filter[T any](arr []T, f func(T) bool) []T {
	var result []T
	for _, v := range arr {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
