package testshlp

import (
	"context"
	"testing"

	"github.com/newity/glog"
	"github.com/stretchr/testify/require"
	"github.com/atomyze-foundation/common-component/loggerhlp"
)

func CreateCtxLogger(t *testing.T) (context.Context, glog.Logger) {
	log, err := loggerhlp.CreateLogger("std", "debug")
	require.NoError(t, err)
	return glog.NewContext(context.Background(), log), log
}

func CreateCtxLoggerWithCancel(t *testing.T) (context.Context, context.CancelFunc, glog.Logger) {
	ctx, log := CreateCtxLogger(t)
	ctx, cancel := context.WithCancel(ctx)
	return ctx, cancel, log
}
