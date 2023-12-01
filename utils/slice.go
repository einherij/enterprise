package utils

func IsInSlice[T comparable](i T, slice []T) bool {
	for _, el := range slice {
		if i == el {
			return true
		}
	}
	return false
}

func First[T any](arr []T) (f T) {
	if len(arr) > 0 {
		return arr[0]
	}
	return f
}

func Last[T any](arr []T) (f T) {
	if len(arr) > 0 {
		return arr[len(arr)-1]
	}
	return f
}

func AreSlicesIntersect[T comparable](slice1, slice2 []T) bool {
	for i := range slice1 {
		for j := range slice2 {
			if slice1[i] == slice2[j] {
				return true
			}
		}
	}
	return false
}

// Range returns array of values [start, end)
func Range[T int | int8 | int16 | int32 | int64 | float32 | float64](start, end T) []T {
	var arr []T
	for i := start; i < end; i++ {
		arr = append(arr, i)
	}
	return arr
}
