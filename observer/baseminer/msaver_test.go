package baseminer

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/atomyze-foundation/common-component/testshlp"
)

func TestSaverRunStop(t *testing.T) {
	ctx, cancel, _ := testshlp.CreateCtxLoggerWithCancel(t)

	runErr := make(chan error)
	go func() {
		runErr <- runSaver[stubCheckPoint, stubPreSaverData, stubSaverData](
			ctx, nil, nil,
			make(chan *saverTask[stubPreSaverData]), 10, nil)
	}()
	cancel()
	require.ErrorIs(t, <-runErr, context.Canceled)
}

func TestSaverNormalWork(t *testing.T) {
	const countTasks = 100
	ctx, cancel, _ := testshlp.CreateCtxLoggerWithCancel(t)

	tasks := make(chan *saverTask[stubPreSaverData])
	stor := newStubStorage(countTasks - 1)

	const psSuff = "_pssuff"
	prsv := newPreSaver(psSuff)

	runErr := make(chan error)
	go func() {
		runErr <- runSaver[stubCheckPoint, stubPreSaverData, stubSaverData](
			ctx, stor, prsv, tasks, 10, nil)
	}()

	for i := 0; i < countTasks; i++ {
		str := make(chan *parseResult[stubPreSaverData], 1)
		st := &saverTask[stubPreSaverData]{
			result: str,
		}
		str <- &parseResult[stubPreSaverData]{
			data: &stubPreSaverData{
				num:   i,
				data1: fmt.Sprintf("data1_%d", i),
				data2: fmt.Sprintf("data2_%d", i),
			},
		}
		tasks <- st
	}

	<-stor.wasTriggered()

	cancel()
	require.ErrorIs(t, <-runErr, context.Canceled)

	require.Len(t, stor.data, countTasks)
	for i := 0; i < countTasks; i++ {
		require.Equal(t, i, stor.data[i].num)
		require.True(t, strings.HasSuffix(stor.data[i].data1, psSuff))
	}
}

func TestSaverWorkWithErrors(t *testing.T) {
	const countTasks, numTaskWithErr = 100, 77
	ctx, cancel, log := testshlp.CreateCtxLoggerWithCancel(t)

	var expectedErrors []error
	addExpectedErrors := func(msg string) error {
		e := errors.New(msg)
		expectedErrors = append(expectedErrors, e)
		return e
	}

	tasks := make(chan *saverTask[stubPreSaverData], countTasks)
	for i := 0; i < countTasks; i++ {
		str := make(chan *parseResult[stubPreSaverData], 1)
		st := &saverTask[stubPreSaverData]{
			result: str,
		}
		if i == numTaskWithErr {
			str <- &parseResult[stubPreSaverData]{
				err: addExpectedErrors(fmt.Sprintf("error task result %v", numTaskWithErr)),
			}
		} else {
			str <- &parseResult[stubPreSaverData]{
				data: &stubPreSaverData{
					num:   i,
					data1: fmt.Sprintf("data1_%d", i),
					data2: fmt.Sprintf("data2_%d", i),
				},
			}
		}

		tasks <- st
	}

	stor := newStubStorage(countTasks - 1)
	stor.callHlp.AddErrMap(stor.Save, map[int]error{
		0: addExpectedErrors("Save error 0"),
		3: addExpectedErrors("Save error 3"),
	})

	const psSuff = "_pssuff"
	prsv := newPreSaver(psSuff)

	go func() {
		<-stor.wasTriggered()
		cancel()
	}()

	var caughtErrors []error
	for {
		err := runSaver[stubCheckPoint, stubPreSaverData, stubSaverData](
			ctx, stor, prsv, tasks, 10, stor.chp)
		if ctx.Err() != nil && errors.Is(err, context.Canceled) {
			break
		}
		log.Errorf("caught error %s", err)
		caughtErrors = append(caughtErrors, err)
	}

	require.ElementsMatch(t, expectedErrors, caughtErrors)
	require.Equal(t, countTasks-1, stor.data[len(stor.data)-1].num)
}
