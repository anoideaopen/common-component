// Package errorshlp contains Newity functions for errors
package errorshlp

import (
	"errors"
	"fmt"
	"io"
)

type (
	// ErrType is a type for error type
	ErrType string
	// ComponentName is a type for component name
	ComponentName string
)

// DetailsError types
type DetailsError struct {
	err       error
	Type      ErrType
	Component ComponentName
}

// Error types
func (de *DetailsError) Error() string {
	return fmt.Sprintf("details component: [%s], src: [%s], error: %s", de.Component, de.Type, de.err)
}

// Unwrap returns the underlying error
func (de *DetailsError) Unwrap() error {
	return de.err
}

// Cause returns the underlying error
func (de *DetailsError) Cause() error {
	return de.err
}

// Format formats the error
func (de *DetailsError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", de.Cause())
			_, _ = io.WriteString(s, de.Error())
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, de.Error())
	}
}

// WrapWithDetails wraps error with details
func WrapWithDetails(err error, errType ErrType, componentName ComponentName) error {
	if err == nil {
		return nil
	}

	// WrapWithDetails wraps only once
	if _, ok := ExtractDetailsError(err); ok {
		return err
	}

	return &DetailsError{
		err:       err,
		Type:      errType,
		Component: componentName,
	}
}

// ExtractDetailsError extracts DetailsError from error
func ExtractDetailsError(err error) (*DetailsError, bool) {
	var de *DetailsError
	if errors.As(err, &de) {
		return de, true
	}
	return nil, false
}
