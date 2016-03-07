// +build !release

package util

import "errors"

var (
	// ErrNotImplimented : not implemented
	ErrNotImplimented = errors.New("not implemented")
)
