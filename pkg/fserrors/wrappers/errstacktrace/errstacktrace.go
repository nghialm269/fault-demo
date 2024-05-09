package errstacktrace

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// withStacktrace implements the error interface and stores a simple table of data.
type withStacktrace struct {
	underlying error
	pcs        []uintptr
}

func (e *withStacktrace) Error() string  { return "<errstacktrace>" }
func (e *withStacktrace) Cause() error   { return e.underlying }
func (e *withStacktrace) Unwrap() error  { return e.underlying }
func (e *withStacktrace) String() string { return e.Error() }

// Implements `StackTrace` method similar to https://github.com/pkg/errors
// This will help Sentry extract the stacktrace: https://github.com/getsentry/sentry-go/blob/23c5137c84c74371307ce682f1283f3c44dd921b/stacktrace.go#L83-L84
func (e *withStacktrace) StackTrace() []uintptr {
	return e.pcs
}

// Wrap wraps an error with stacktrace.
func Wrap(err error, skip int) error {
	if err == nil {
		return nil
	}

	return &withStacktrace{err, callers(skip + 1)}
}

// With implements the Fault Wrapper interface.
func With(skip int) func(error) error {
	return func(err error) error {
		return Wrap(err, skip+1)
	}
}

// Unwrap returns the first stacktrace in the error chain
// or `nil` if no stacktrace found.
func Unwrap(err error) []uintptr {
	for err != nil {
		if f, ok := err.(*withStacktrace); ok {
			return f.pcs
		}

		err = errors.Unwrap(err)
	}

	return nil
}

// Get returns the first stacktrace in the error chain as string
// or empty string if no stacktrace found.
func Get(err error) string {
	for err != nil {
		if f, ok := err.(*withStacktrace); ok {
			return formatStacktrace(f.pcs)
		}

		err = errors.Unwrap(err)
	}

	return ""
}

func formatStacktrace(pcs []uintptr) string {
	frames := runtime.CallersFrames(pcs)

	var sb strings.Builder

	for {
		frame, more := frames.Next()

		function := frame.Function
		if function == "" {
			function = "<unknown function>"
		}

		file := frame.File
		if file == "" {
			file = "<unknown file>"
		}

		sb.WriteString(fmt.Sprintf("%s %s:%d", function, file, frame.Line))

		if !more {
			break
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

func HasStacktrace(err error, skip int) bool {
	var stackTraceErr *withStacktrace
	if errors.As(err, &stackTraceErr) {
		return ancestorOfCause(callers(skip+1), stackTraceErr.pcs)
	}

	return false
}

// Code below are from https://gist.github.com/lawrencejones/3a392f7116220a9799e55460fa57622d
// Read more: https://incident.io/blog/golang-errors

// ancestorOfCause returns true if the caller looks to be an ancestor of the given stack
// trace. We check this by seeing whether our stack prefix-matches the cause stack, which
// should imply the error was generated directly from our goroutine.
func ancestorOfCause(ourStack []uintptr, causeStack []uintptr) bool {
	// Stack traces are ordered such that the deepest frame is first. We'll want to check
	// for prefix matching in reverse.
	//
	// As an example, imagine we have a prefix-matching stack for ourselves:
	// [
	//   "github.com/onsi/ginkgo/internal/leafnodes.(*runner).runSync",
	//   "github.com/incident-io/core/server/pkg/errors_test.TestSuite",
	//   "testing.tRunner",
	//   "runtime.goexit"
	// ]
	//
	// We'll want to compare this against an error cause that will have happened further
	// down the stack. An example stack trace from such an error might be:
	// [
	//   "github.com/incident-io/core/server/pkg/errors.New",
	//   "github.com/incident-io/core/server/pkg/errors_test.glob..func1.2.2.2.1",,
	//   "github.com/onsi/ginkgo/internal/leafnodes.(*runner).runSync",
	//   "github.com/incident-io/core/server/pkg/errors_test.TestSuite",
	//   "testing.tRunner",
	//   "runtime.goexit"
	// ]
	//
	// They prefix match, but we'll have to handle the match carefully as we need to match
	// from back to forward.

	// We can't possibly prefix match if our stack is larger than the cause stack.
	if len(ourStack) > len(causeStack) {
		return false
	}

	// We know the sizes are compatible, so compare program counters from back to front.
	for idx := 0; idx < len(ourStack); idx++ {
		if ourStack[len(ourStack)-1-idx] != (uintptr)(causeStack[len(causeStack)-1-idx]) {
			return false
		}
	}

	// All comparisons checked out, these stacks match
	return true
}

func callers(skip int) []uintptr {
	pc := make([]uintptr, 32)        // assume we'll have at most 32 frames
	n := runtime.Callers(skip+3, pc) // capture those frames, skipping runtime.Callers, ourself and the calling function

	return pc[:n] // return everything that we captured
}
