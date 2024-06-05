package baseminer

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrNilTask is an error for nil task
	ErrNilTask = errors.New("nil task")
	// ErrParserTasksWasClosed is an error for closed tasks
	ErrParserTasksWasClosed = errors.New("tasks was closed")
)

func runParserPool[TSrcData any, TPreSaverData any](
	ctx context.Context,
	prsr DataParser[TSrcData, TPreSaverData],
	countWorkers uint,
	tasks <-chan *parserTask[TSrcData, TPreSaverData],
) error {
	eg, ctxEg := errgroup.WithContext(ctx)

	for range countWorkers {
		eg.Go(func() error {
			return parserWorker(ctxEg, prsr, tasks)
		})
	}
	return eg.Wait()
}

func parserWorker[TSrcData any, TPreSaverData any](
	ctx context.Context,
	prsr DataParser[TSrcData, TPreSaverData],
	tasks <-chan *parserTask[TSrcData, TPreSaverData],
) error {
	tasksL := tasks
	var res *parseResult[TPreSaverData]
	var taskResult chan<- *parseResult[TPreSaverData]

	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case task, ok := <-tasksL:
			if !ok {
				return ErrParserTasksWasClosed
			}
			if task == nil {
				return ErrNilTask
			}
			d, err := prsr.Parse(ctx, task.srcData)
			res = &parseResult[TPreSaverData]{
				err:  err,
				data: d,
			}
			tasksL = nil
			taskResult = task.result
		case taskResult <- res:
			taskResult = nil
			tasksL = tasks
		}
	}
	return ctx.Err()
}
