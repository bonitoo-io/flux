package flux

import (
	"strings"

	"github.com/influxdata/flux/codes"
)

// Error is the error struct of flux.
type Error struct {
	// Code is the code of the error as defined in the codes package.
	// This describes the type and category of the error. It is required.
	Code codes.Code

	// Msg contains a human-readable description and additional information
	// about the error itself. This is optional.
	Msg string

	// Err contains the error that was the cause of this error.
	// This is optional.
	Err error
}

// Error implement the error interface by outputting the Code and Err.
func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString(e.Code.String())
	if e.Msg != "" {
		b.WriteString(": ")
		b.WriteString(e.Msg)
	}
	if e.Err != nil {
		b.WriteString(": ")
		b.WriteString(e.Err.Error())
	}
	return b.String()
}

// Unwrap will return the wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}
