package baseminer

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessorSendNilSrcData(t *testing.T) {
	srcCh := make(chan *stubSrcData)
	parserTasks := make(chan *parserTask[stubSrcData, stubPreSaverData])
	saverTasks := make(chan *saverTask[stubPreSaverData])

	go func() {
		srcCh <- nil
	}()
	require.ErrorIs(t, runProcessor[stubSrcData, stubPreSaverData](
		context.Background(), srcCh, parserTasks, saverTasks), ErrSrcDataIsNil)
}

func TestProcessorNormalWork(t *testing.T) {
	const srcDataCount = 2000

	srcCh := make(chan *stubSrcData, srcDataCount)
	for i := 0; i < srcDataCount; i++ {
		srcCh <- &stubSrcData{num: i + 1}
	}

	parserTasks := make(chan *parserTask[stubSrcData, stubPreSaverData], srcDataCount)
	saverTasks := make(chan *saverTask[stubPreSaverData], srcDataCount)

	allRead := make(chan struct{})
	ptCount, stCount := 0, 0
	go func() {
		defer close(allRead)

		lastBnPt := 0
		for ptCount < srcDataCount {
			parserTask, ok := <-parserTasks
			if !ok {
				return
			}
			ptCount++

			if lastBnPt >= parserTask.srcData.num {
				log.Printf("unordered pipe: got src number %d, but previous src num was %d", parserTask.srcData.num, lastBnPt)
				return
			}
			lastBnPt = parserTask.srcData.num
			parserTask.result <- &parseResult[stubPreSaverData]{data: &stubPreSaverData{num: parserTask.srcData.num}}
		}

		lastBnSt := 0
		for stCount < srcDataCount {
			saverTask, ok := <-saverTasks
			if !ok {
				return
			}
			stCount++

			res := <-saverTask.result
			if lastBnSt >= res.data.num {
				log.Printf("unordered pipe: got src number %d, but previous num was %d", res.data.num, lastBnSt)
				return
			}
			lastBnSt = res.data.num
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	runProcessorErr := make(chan error)
	go func() {
		runProcessorErr <- runProcessor(ctx, srcCh, parserTasks, saverTasks)
	}()

	<-allRead

	cancel()
	require.ErrorIs(t, <-runProcessorErr, context.Canceled)

	require.Equal(t, srcDataCount, ptCount)
	require.Equal(t, srcDataCount, stCount)
}

func TestProcessorAfterCloseSrcDataErr(t *testing.T) {
	srcCh := make(chan *stubSrcData)
	parserTasks := make(chan *parserTask[stubSrcData, stubPreSaverData])
	saverTasks := make(chan *saverTask[stubPreSaverData])

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runProcessorErr := make(chan error)
	go func() {
		runProcessorErr <- runProcessor[stubSrcData, stubPreSaverData](ctx, srcCh, parserTasks, saverTasks)
	}()

	close(srcCh)

	require.ErrorIs(t, <-runProcessorErr, ErrSrcDataClosed)
}
