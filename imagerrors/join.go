// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style

package imagerrors

import (
	"strings"
)

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

type joinError struct {
	errs []error
}

func (e *joinError) Error() string {
	// Since Join returns nil if every value in errs is nil,
	// e.errs cannot be empty.
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	var builder strings.Builder
	builder.WriteString(e.errs[0].Error())
	for _, err := range e.errs[1:] {
		builder.WriteString("\n")
		builder.WriteString(err.Error())
	}
	return builder.String()
}

func (e *joinError) Unwrap() []error {
	return e.errs
}
