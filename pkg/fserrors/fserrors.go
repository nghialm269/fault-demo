package fserrors

import (
	"github.com/Southclaws/fault"

	"github.com/nghialm269/fault-demo/pkg/fserrors/wrappers/errstacktrace"
)

type Wrapper = fault.Wrapper

// NewSentinel is like `New` but without stacktrace.
func NewSentinel(message string, w ...Wrapper) error {
	return fault.New(message, w...)
}

func New(message string, w ...Wrapper) error {
	w = append([]Wrapper{errstacktrace.With(1)}, w...)
	return fault.New(message, w...)
}

func Wrap(err error, w ...Wrapper) error {
	if !errstacktrace.HasStacktrace(err, 1) {
		w = append([]Wrapper{errstacktrace.With(1)}, w...)
	}
	return fault.Wrap(err, w...)
}
