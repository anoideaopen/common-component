package baseminer

import (
	"context"
	"testing"

	"github.com/atomyze-foundation/common-component/testshlp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestParserPoolStartStop(t *testing.T) {
	ctx, cancel, _ := testshlp.CreateCtxLoggerWithCancel(t)

	prsr := newStubParser(nil)

	runErr := make(chan error)
	go func() {
		runErr <- runParserPool[stubSrcData, stubPreSaverData](
			ctx, prsr, 10, make(chan *parserTask[stubSrcData, stubPreSaverData]))
	}()
	cancel()

	require.ErrorIs(t, <-runErr, context.Canceled)
}

func TestParserPoolNormalWork(t *testing.T) {
	const tCount, wCount = 1000, 10

	pTasks := make(chan *parserTask[stubSrcData, stubPreSaverData], tCount)

	// push tasks to channel
	var pTasksResults []chan *parseResult[stubPreSaverData]
	for i := 0; i < tCount; i++ {
		resultch := make(chan *parseResult[stubPreSaverData], 1)

		pTasksResults = append(pTasksResults, resultch)
		pTasks <- &parserTask[stubSrcData, stubPreSaverData]{
			srcData: &stubSrcData{num: i},
			result:  resultch,
		}
	}

	runErr := make(chan error)
	go func() {
		runErr <- runParserPool[stubSrcData, stubPreSaverData](
			context.Background(), newStubParser(nil), wCount, pTasks)
	}()

	// check results
	for i, rch := range pTasksResults {
		r, ok := <-rch
		require.True(t, ok)
		require.NoError(t, r.err)

		require.Equal(t, i, r.data.num)
	}

	// pool still work
	select {
	case err := <-runErr:
		require.FailNow(t, "runParserPool finished", err)
	default:
	}
}

func TestParserPoolWorkWithParserErrors(t *testing.T) {
	const tCount, wCount = 1000, 10

	pTasks := make(chan *parserTask[stubSrcData, stubPreSaverData], tCount)

	// push tasks to channel
	var pTasksResults []chan *parseResult[stubPreSaverData]
	for i := 0; i < tCount; i++ {
		resultch := make(chan *parseResult[stubPreSaverData], 1)

		pTasksResults = append(pTasksResults, resultch)
		pTasks <- &parserTask[stubSrcData, stubPreSaverData]{
			srcData: &stubSrcData{num: i},
			result:  resultch,
		}
	}

	badSrcData := map[int]error{
		12: errors.New("data 12"),
		33: errors.New("data 33"),
	}

	runErr := make(chan error)
	go func() {
		runErr <- runParserPool[stubSrcData, stubPreSaverData](
			context.Background(),
			newStubParser(badSrcData), wCount, pTasks)
	}()

	// check results
	for i, rch := range pTasksResults {
		r, ok := <-rch
		require.True(t, ok)
		if err, ok := badSrcData[i]; ok {
			require.ErrorIs(t, r.err, err)
		} else {
			require.NoError(t, r.err)
			require.Equal(t, i, r.data.num)
		}
	}

	// pool still work
	select {
	case err := <-runErr:
		require.FailNow(t, "runParserPool finished", err)
	default:
	}
}

func TestParserPoolAfterCloseTasks(t *testing.T) {
	const tCount, wCount = 1000, 10

	pTasks := make(chan *parserTask[stubSrcData, stubPreSaverData], tCount)
	for i := 0; i < tCount; i++ {
		pTasks <- &parserTask[stubSrcData, stubPreSaverData]{
			srcData: &stubSrcData{num: i},
			result:  make(chan *parseResult[stubPreSaverData], 1),
		}
	}
	close(pTasks)

	ctx, _, _ := testshlp.CreateCtxLoggerWithCancel(t)

	prsr := newStubParser(nil)

	err := runParserPool[stubSrcData, stubPreSaverData](ctx, prsr, wCount, pTasks)
	require.ErrorIs(t, err, ErrParserTasksWasClosed)
}
