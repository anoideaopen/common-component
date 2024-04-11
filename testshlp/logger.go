package testshlp

import (
	"context"
	"testing"

	"github.com/anoideaopen/common-component/loggerhlp"
	"github.com/anoideaopen/glog"
	"github.com/stretchr/testify/require"
)

// CreateCtxLogger creates a context with logger
func CreateCtxLogger(t *testing.T) (context.Context, glog.Logger) {
	log, err := loggerhlp.CreateLogger("std", "debug")
	require.NoError(t, err)
	return glog.NewContext(context.Background(), log), log
}

// CreateCtxLoggerWithCancel creates a context with logger and cancel function
func CreateCtxLoggerWithCancel(t *testing.T) (context.Context, context.CancelFunc, glog.Logger) {
	ctx, log := CreateCtxLogger(t)
	ctx, cancel := context.WithCancel(ctx)
	return ctx, cancel, log
}
