package specification

// Specification is the base interface for all specifications in DDD.
// A specification encapsulates a business rule that can be used to:
// - Validate entities
// - Query/filter entities
// - Check business conditions
//
// The generic type T represents the entity type that the specification operates on.
type Specification[T any] interface {
	// IsSatisfiedBy checks if the given entity satisfies this specification
	IsSatisfiedBy(entity T) bool
}

// BuilderSpecification wraps a Specification to provide builder pattern methods (And, Or, Not)
type BuilderSpecification[T any] struct {
	spec Specification[T]
}

// NewBuilder wraps a specification to provide builder API
func NewBuilder[T any](spec Specification[T]) *BuilderSpecification[T] {
	return &BuilderSpecification[T]{spec: spec}
}

// IsSatisfiedBy delegates to the wrapped specification
func (b *BuilderSpecification[T]) IsSatisfiedBy(entity T) bool {
	return b.spec.IsSatisfiedBy(entity)
}

// And combines this specification with another using logical AND
func (b *BuilderSpecification[T]) And(other Specification[T]) *BuilderSpecification[T] {
	return NewBuilder(NewAndSpecification(b.spec, other))
}

// Or combines this specification with another using logical OR
func (b *BuilderSpecification[T]) Or(other Specification[T]) *BuilderSpecification[T] {
	return NewBuilder(NewOrSpecification(b.spec, other))
}

// Not negates this specification using logical NOT
func (b *BuilderSpecification[T]) Not() *BuilderSpecification[T] {
	return NewBuilder(NewNotSpecification(b.spec))
}

// Spec returns the underlying specification
func (b *BuilderSpecification[T]) Spec() Specification[T] {
	return b.spec
}
