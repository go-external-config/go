package collection

func SliceToSet[T comparable](slice []T) map[T]any {
	result := make(map[T]any)
	for _, value := range slice {
		result[value] = nil
	}
	return result
}

// a - b
func SubtractSlice[T comparable](a, b []T) []T {
	set := SliceToSet(b)
	var diff []T
	for _, x := range a {
		if _, found := set[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
