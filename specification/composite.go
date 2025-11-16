package specification

// AndSpecification combines two specifications with logical AND
type AndSpecification[T any] struct {
	left  Specification[T]
	right Specification[T]
}

// NewAndSpecification creates a new AND specification
func NewAndSpecification[T any](left, right Specification[T]) Specification[T] {
	return &AndSpecification[T]{
		left:  left,
		right: right,
	}
}

func (s *AndSpecification[T]) IsSatisfiedBy(entity T) bool {
	return s.left.IsSatisfiedBy(entity) && s.right.IsSatisfiedBy(entity)
}

// OrSpecification combines two specifications with logical OR
type OrSpecification[T any] struct {
	left  Specification[T]
	right Specification[T]
}

// NewOrSpecification creates a new OR specification
func NewOrSpecification[T any](left, right Specification[T]) Specification[T] {
	return &OrSpecification[T]{
		left:  left,
		right: right,
	}
}

func (s *OrSpecification[T]) IsSatisfiedBy(entity T) bool {
	return s.left.IsSatisfiedBy(entity) || s.right.IsSatisfiedBy(entity)
}

// NotSpecification negates a specification with logical NOT
type NotSpecification[T any] struct {
	spec Specification[T]
}

// NewNotSpecification creates a new NOT specification
func NewNotSpecification[T any](spec Specification[T]) Specification[T] {
	return &NotSpecification[T]{
		spec: spec,
	}
}

func (s *NotSpecification[T]) IsSatisfiedBy(entity T) bool {
	return !s.spec.IsSatisfiedBy(entity)
}

