package utils

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
