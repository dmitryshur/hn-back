package data

func ToPointer[T any](n T) *T {
	return &n
}

func IfElse[T any](condition bool, true T, false T) T {
	if condition {
		return true
	}

	return false
}
