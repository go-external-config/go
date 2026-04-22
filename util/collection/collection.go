package collection

import (
	"cmp"
	"slices"
)

func SliceToSet[T comparable](slice []T) map[T]any {
	result := make(map[T]any)
	for _, value := range slice {
		result[value] = nil
	}
	return result
}

// new slice a - b
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

// new reverced slice
func ReverseSlice[T any](slice []T) []T {
	reversedSlice := make([]T, len(slice))
	copy(reversedSlice, slice)
	slices.Reverse(reversedSlice)
	return reversedSlice
}

// the same sorted slice
func Sort[T cmp.Ordered](slice []T) []T {
	slices.Sort(slice)
	return slice
}

func Distinct[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	out := make([]T, 0, len(slice))

	for _, v := range slice {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
