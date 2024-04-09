package baseminer

import (
	"context"
	"errors"
	"testing"

	"github.com/anoideaopen/common-component/testshlp"
	"github.com/stretchr/testify/require"
)

func TestHandlerStartStop(t *testing.T) {
	ctx, cancel, _ := testshlp.CreateCtxLoggerWithCancel(t)

	stor := newStubStorage(0)
	prsr := newStubParser(nil)

	createSrc := func(ctx context.Context, startFromChp *stubCheckPoint) (DataSource[stubSrcData], error) {
		return newStubSrc(ctx, nil, startFromChp), nil
	}

	createPreSaver := func(ctx context.Context) (PreSaver[stubCheckPoint, stubPreSaverData, stubSaverData], error) {
		return newPreSaver("_ps"), nil
	}

	hndlr := NewMHandler[stubCheckPoint, stubSrcData, stubPreSaverData, stubSaverData](
		ctx, stor, createSrc, prsr, createPreSaver,
		1, 1, 1)

	require.NotNil(t, hndlr)

	runFinishErr := make(chan error)
	go func() {
		runFinishErr <- hndlr.RunHandler(ctx)
	}()

	cancel()
	err := <-runFinishErr

	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestHandlerNormalWork(t *testing.T) {
	const (
		parserPoolCountWorkers = 10
		saverTasksBufSize      = 10
		saverMaxBatchSize      = 10
	)
	ctx, cancel, _ := testshlp.CreateCtxLoggerWithCancel(t)

	allSrcData := createSrcData(0, 1000)
	stor := newStubStorage(allSrcData[len(allSrcData)-1].num)
	prsr := newStubParser(nil)

	srcDataCr := &stubCrSrc{
		srcData: allSrcData,
	}

	createPreSaver := func(ctx context.Context) (PreSaver[stubCheckPoint, stubPreSaverData, stubSaverData], error) {
		return newPreSaver("_ps"), nil
	}

	hndlr := NewMHandler[stubCheckPoint, stubSrcData, stubPreSaverData, stubSaverData](
		ctx, stor, srcDataCr.create, prsr, createPreSaver,
		parserPoolCountWorkers, saverTasksBufSize, saverMaxBatchSize)

	require.NotNil(t, hndlr)

	runFinishErr := make(chan error)
	go func() {
		runFinishErr <- hndlr.RunHandler(ctx)
	}()

	<-stor.wasTriggered()
	cancel()
	err := <-runFinishErr

	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)

	require.Equal(t, len(allSrcData), len(stor.data))
	for i := 0; i < len(allSrcData); i++ {
		require.Equal(t, allSrcData[i].num, stor.data[i].num)
	}
}

func TestHandlerWorkWithErrors(t *testing.T) {
	const (
		parserPoolCountWorkers = 10
		saverTasksBufSize      = 10
		saverMaxBatchSize      = 10
	)

	ctx, cancel, log := testshlp.CreateCtxLoggerWithCancel(t)

	allSrcData := createSrcData(0, 1000)

	var expectedErrors []error
	addExpectedErrors := func(msg string) error {
		e := errors.New(msg)
		expectedErrors = append(expectedErrors, e)
		return e
	}

	stor := newStubStorage(allSrcData[len(allSrcData)-1].num)
	stor.callHlp.AddErrMap(stor.FindCheckPoint, map[int]error{
		0: addExpectedErrors("FindCheckPoint error 0"),
		3: addExpectedErrors("FindCheckPoint error 3"),
	})
	stor.callHlp.AddErrMap(stor.Save, map[int]error{
		0:  addExpectedErrors("Save error 0"),
		10: addExpectedErrors("Save error 10"),
	})

	srcDataCr := &stubCrSrc{
		srcData: allSrcData,
	}
	srcDataCr.callHlp.AddErrMap(srcDataCr.create, map[int]error{
		0: addExpectedErrors("create src error 0"),
		2: addExpectedErrors("create src error 2"),
	})

	prsr := newStubParser(nil)

	createPreSaver := func(ctx context.Context) (PreSaver[stubCheckPoint, stubPreSaverData, stubSaverData], error) {
		return newPreSaver("_ps"), nil
	}

	hndlr := NewMHandler[stubCheckPoint, stubSrcData, stubPreSaverData, stubSaverData](
		ctx, stor, srcDataCr.create, prsr, createPreSaver,
		parserPoolCountWorkers, saverTasksBufSize, saverMaxBatchSize)
	require.NotNil(t, hndlr)

	go func() {
		<-stor.wasTriggered()
		cancel()
	}()

	var caughtErrors []error
	for {
		err := hndlr.RunHandler(ctx)
		if ctx.Err() != nil && errors.Is(err, context.Canceled) {
			break
		}
		log.Errorf("caught error %s", err)
		caughtErrors = append(caughtErrors, err)
	}
	require.ElementsMatch(t, expectedErrors, caughtErrors)

	require.Equal(t, len(allSrcData), len(stor.data))
	for i := 0; i < len(allSrcData); i++ {
		require.Equal(t, allSrcData[i].num, stor.data[i].num)
	}
}

func TestHandlerWorkWithPermanentBadBlock(t *testing.T) {
	const (
		parserPoolCountWorkers = 10
		saverTasksBufSize      = 10
		saverMaxBatchSize      = 10
		errBlockNum            = 500
	)
	ctx, _ := testshlp.CreateCtxLogger(t)

	allSrcData := createSrcData(0, 1000)

	stor := newStubStorage(allSrcData[len(allSrcData)-1].num)

	srcDataCr := &stubCrSrc{
		srcData: allSrcData,
	}

	createPreSaver := func(ctx context.Context) (PreSaver[stubCheckPoint, stubPreSaverData, stubSaverData], error) {
		return newPreSaver("_ps"), nil
	}

	blErr := errors.New("block 501 error")
	prsr := newStubParser(map[int]error{errBlockNum: blErr})

	hndlr := NewMHandler[stubCheckPoint, stubSrcData, stubPreSaverData, stubSaverData](
		ctx, stor, srcDataCr.create, prsr, createPreSaver,
		parserPoolCountWorkers, saverTasksBufSize, saverMaxBatchSize)

	require.NotNil(t, hndlr)

	// run after error will produce same result
	for i := 0; i < 2; i++ {
		err := hndlr.RunHandler(ctx)
		require.ErrorIs(t, err, blErr)
		require.Equal(t, errBlockNum, len(stor.data))
	}
}

func createSrcData(firstNum int, count int) []*stubSrcData {
	res := make([]*stubSrcData, 0, count)
	for i := 0; i < count; i++ {
		res = append(res, &stubSrcData{num: firstNum + i})
	}
	return res
}
