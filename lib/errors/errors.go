package errors

import (
	"github.com/palantir/stacktrace"
)

// Code - Max value is 4294967295
type Code = stacktrace.ErrorCode

func init() {
	stacktrace.DefaultFormat = stacktrace.FormatFull
}

// ErrCode extracts the error code from an error.
var ErrCode = stacktrace.GetCode

// New is a drop-in replacement for fmt.Errorf that includes line number information.
var New = stacktrace.NewError

// NewWithCode is similar to New but also attaches an error code.
var NewWithCode = stacktrace.NewErrorWithCode

// Wrap an error to include line number information.
var Wrap = stacktrace.Propagate

// WrapWithCode is similar to Wrap but also attaches an error code.
var WrapWithCode = stacktrace.PropagateWithCode

// Wrapf is similar to Wrap but the msg and vals arguments work like the ones for fmt.Errorf.
var Wrapf = stacktrace.Propagate
