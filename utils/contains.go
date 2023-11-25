package utils

func IsInSlice[T comparable](i T, slice []T) bool {
	for _, el := range slice {
		if i == el {
			return true
		}
	}
	return false
}
