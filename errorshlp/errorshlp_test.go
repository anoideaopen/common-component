package errorshlp

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const errTypeTest ErrType = "test"

func TestWrapWithDetails(t *testing.T) {
	err := errors.New("someTestErr")
	err = WrapWithDetails(err, errTypeTest, "Test1")
	err = WrapWithDetails(err, errTypeTest, "Test2")

	errDetails, ok := ExtractDetailsError(err)
	require.True(t, ok)
	require.NotNil(t, errDetails)
	require.EqualValues(t, errDetails.Component, "Test1")

	// fmt.Printf("%+v", err)
}

func TestWrapWithDetailsStackTraceDoesntOverride(t *testing.T) {
	err := errors.New("someTestErr")
	origErr := err
	err = WrapWithDetails(origErr, errTypeTest, "Test")
	err = WrapWithDetails(err, errTypeTest, "Test2")
	err = errors.WithStack(err)

	// fmt.Println("err *********************")
	// fmt.Printf("%+v\n", err)
	// fmt.Println("origErr *********************")
	// fmt.Printf("%+v\n", origErr)
	// fmt.Println("*********************")

	st1 := extractStackTrace(origErr)
	st2 := extractStackTrace(err)
	require.EqualValues(t, len(st1), len(st2))
}

func TestWrapWithDetailsStackTraceFromDetailsError(t *testing.T) {
	origErr := errors.New("someTestErr")
	err := WrapWithDetails(origErr, errTypeTest, "Test")
	err = WrapWithDetails(err, errTypeTest, "Test2")

	dErr, ok := ExtractDetailsError(err)
	require.True(t, ok)

	st1 := extractStackTrace(origErr)
	st2 := extractStackTrace(dErr.Cause())
	require.EqualValues(t, len(st1), len(st2))
}

func TestExtractDetailsError(t *testing.T) {
	err := errors.New("someTestErr")
	err = WrapWithDetails(err, errTypeTest, "Test")
	err = errors.Wrap(err, "additional wrap1")
	err = errors.Wrap(err, "additional wrap2")

	errDetails, ok := ExtractDetailsError(err)
	require.True(t, ok)
	require.NotNil(t, errDetails)
	require.EqualValues(t, errDetails.Type, errTypeTest)
	require.EqualValues(t, errDetails.Component, "Test")

	// fmt.Printf("%+v", err)
}

func TestWrapWithDetailsErrNil(t *testing.T) {
	errNil := WrapWithDetails(nil, errTypeTest, "Test")
	require.Nil(t, errNil)
}

func extractStackTrace(err error) errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	st, ok := err.(stackTracer) //nolint:errorlint
	if !ok {
		return nil
	}
	return st.StackTrace()
}
