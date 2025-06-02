package typ

import "iter"

// Filter2 is a function that filters a iter.Seq2 based on a predicate.
func Filter2[T any](seq iter.Seq2[T, error], predicate func(T) bool) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		seq(func(item T, err error) bool {
			if err != nil {
				return yield(item, err)
			}
			if predicate(item) {
				return yield(item, nil)
			}
			return true // continue iterating
		})
	}
}

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
