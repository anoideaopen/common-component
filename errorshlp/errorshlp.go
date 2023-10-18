// Package errorshlp contains Newity functions for errors
package errorshlp

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type (
	ErrType       string
	ComponentName string
)

type DetailsError struct {
	err       error
	Type      ErrType
	Component ComponentName
}

func (de *DetailsError) Error() string {
	return fmt.Sprintf("details component: [%s], src: [%s], error: %s", de.Component, de.Type, de.err)
}

func (de *DetailsError) Unwrap() error {
	return de.err
}

func (de *DetailsError) Cause() error {
	return de.err
}

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

func ExtractDetailsError(err error) (*DetailsError, bool) {
	var de *DetailsError
	if errors.As(err, &de) {
		return de, true
	}
	return nil, false
}
