package data

func ToPointer[T any](n T) *T {
	return &n
}
