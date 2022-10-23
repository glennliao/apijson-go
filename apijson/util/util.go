package util

func Contains[T comparable](arr []T, item T) bool {
	for _, t := range arr {
		if t == item {
			return true
		}
	}
	return false
}

func Reverse[T any](arr *[]T) {
	length := len(*arr)
	for i := 0; i < length/2; i++ {
		(*arr)[i], (*arr)[length-1-i] = (*arr)[length-1-i], (*arr)[i]
	}
}
