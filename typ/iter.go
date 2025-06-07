package typ

import "iter"

// Filter is a function that filters a iter.Seq based on a predicate.
func Filter[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		seq(func(item T) bool {
			if predicate(item) {
				return yield(item)
			}
			return true // continue iterating
		})
	}
}

// Filter2 is a function that filters a iter.Seq based on a predicate.
func Filter2[T any, U any](seq iter.Seq2[T, U], predicate func(T, U) bool) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		seq(func(t T, u U) bool {
			if predicate(t, u) {
				return yield(t, u)
			}
			return true // continue iterating
		})
	}
}

// FilterE is a function that filters a iter.Seq2 with an error based on a predicate.
func FilterE[T any](seq iter.Seq2[T, error], predicate func(T) (bool, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		seq(func(item T, err error) bool {
			if err != nil {
				return yield(item, err)
			}
			m, err := predicate(item)
			if err != nil {
				return yield(item, err)
			}
			if m {
				return yield(item, nil)
			}
			return true // continue iterating
		})
	}
}

// WithNilError is a function that wraps a iter.Seq and yields items with a nil error.
func WithNilError[T any](seq iter.Seq[T]) iter.Seq2[T, error] {
	return WithError(seq, nil)
}

// WithError is a function that wraps a iter.Seq and yields items with a error.
func WithError[T any](seq iter.Seq[T], err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		seq(func(item T) bool {
			return yield(item, err)
		})
	}
}

// CollectWithError is a function that collects items from a iter.Seq2 with an error into a slice, returning an error if any.
func CollectWithError[T any](seq iter.Seq2[T, error]) (res []T, _ error) {
	for item, err := range seq {
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}
