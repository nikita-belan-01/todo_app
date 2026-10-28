package domain

type Nullable[T any] struct {
	Value *T
	Set   bool
}

func applyNotNull[T comparable](dst *T, n Nullable[T]) bool {
	if !n.Set || n.Value == nil || *n.Value == *dst {
		return false
	}

	*dst = *n.Value

	return true
}

func applyNullable[T comparable](dst **T, n Nullable[T]) bool {
	if !n.Set || equalPtr(n.Value, *dst) {
		return false
	}

	*dst = n.Value

	return true
}

func equalPtr[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}

	return *a == *b
}
