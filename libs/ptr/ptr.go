package ptr

func Ptr[T any](v T) *T {
	return &v
}

func DePtr[T any](v *T) T {
	return *v
}
