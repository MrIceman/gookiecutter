package collectionoperation

func FlattenMapValues[T comparable, U any](m map[T][]U) []U {
	var result []U
	for _, v := range m {
		result = append(result, v...)
	}
	return result
}
