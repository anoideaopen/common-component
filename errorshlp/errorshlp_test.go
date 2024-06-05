package errorshlp

import (
	"errors"
	"fmt"
	"testing"

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

func TestExtractDetailsError(t *testing.T) {
	err := errors.New("someTestErr")
	err = WrapWithDetails(err, errTypeTest, "Test")
	err = fmt.Errorf("additional wrap1: %w", err)
	err = fmt.Errorf("additional wrap2: %w", err)

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
