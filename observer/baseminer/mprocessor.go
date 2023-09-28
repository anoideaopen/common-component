package baseminer

import (
	"context"

	"github.com/pkg/errors"
)

var (
	// ErrSrcDataClosed is an error for closed src data
	ErrSrcDataClosed = errors.New("source of src data was closed")
	// ErrSrcDataIsNil is an error for nil src data
	ErrSrcDataIsNil = errors.New("src data is nil")
)

func runProcessor[TSrcData any, TPreSaverData any](ctx context.Context,
	src <-chan *TSrcData,
	parserTasks chan<- *parserTask[TSrcData, TPreSaverData],
	saverTasks chan<- *saverTask[TPreSaverData],
) error {
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case srcData, ok := <-src:
			if !ok {
				return errors.WithStack(ErrSrcDataClosed)
			}

			if srcData == nil {
				return errors.WithStack(ErrSrcDataIsNil)
			}

			if err := queueToProcess(ctx, srcData, parserTasks, saverTasks); err != nil {
				return err
			}
		}
	}

	return ctx.Err()
}

func queueToProcess[TSrcData any, TPreSaverData any](
	ctx context.Context, srcData *TSrcData,
	parserTasks chan<- *parserTask[TSrcData, TPreSaverData],
	saverTasks chan<- *saverTask[TPreSaverData],
) error {
	resultCh := make(chan *parseResult[TPreSaverData], 1)

	saverTask := &saverTask[TPreSaverData]{
		result: resultCh,
	}

	parserTask := &parserTask[TSrcData, TPreSaverData]{
		srcData: srcData,
		result:  resultCh,
	}

	saverL, parserL := saverTasks, parserTasks

	for saverL != nil || parserL != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case saverL <- saverTask:
			saverL = nil
		case parserL <- parserTask:
			parserL = nil
		}
	}

	return nil
}
