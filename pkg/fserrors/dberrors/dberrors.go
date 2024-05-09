package dberrors

import "github.com/nghialm269/fault-demo/pkg/fserrors"

var (
	ErrEntryNotFound  = fserrors.NewSentinel("entry not found")
	ErrDuplicateEntry = fserrors.NewSentinel("duplicate entry")
)
