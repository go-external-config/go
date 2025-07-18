package lang

func If[T any](cond bool, v1, v2 T) T {
	if cond {
		return v1
	}
	return v2
}

func FirstNonEmpty[T comparable](values ...T) T {
	empty := *new(T)
	for _, value := range values {
		if value != empty {
			return value
		}
	}
	return empty
}
