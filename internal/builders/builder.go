package builders

type Builder[T interface{}] interface {
	Build() T
}
