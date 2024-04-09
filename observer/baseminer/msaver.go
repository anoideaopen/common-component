package baseminer

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// ErrSaverTasksWasClosed is an error for closed tasks
var ErrSaverTasksWasClosed = errors.New("tasks was closed")

func runSaver[TCheckPoint any, TPreSaverData any, TSaverData any](
	ctx context.Context,
	stor Storage[TCheckPoint, TSaverData],
	prsv PreSaver[TCheckPoint, TPreSaverData, TSaverData],
	tasks <-chan *saverTask[TPreSaverData],
	maxBatchSize uint, initCheckPoint *TCheckPoint,
) error {
	eg, egCtx := errgroup.WithContext(ctx)

	results := make(chan *parseResult[TPreSaverData], maxBatchSize)
	defer close(results)

	eg.Go(func() error {
		return poolTasksResults(egCtx, tasks, results)
	})

	eg.Go(func() error {
		return saveLoop(egCtx, stor, prsv, results, maxBatchSize, initCheckPoint)
	})

	return eg.Wait()
}

func poolTasksResults[TPreSaverData any](ctx context.Context,
	tasks <-chan *saverTask[TPreSaverData],
	results chan<- *parseResult[TPreSaverData],
) error {
	tasksL := tasks
	var taskResL <-chan *parseResult[TPreSaverData]
	var resultsL chan<- *parseResult[TPreSaverData]
	var result *parseResult[TPreSaverData]
	for ctx.Err() == nil {
		select {
		case task, ok := <-tasksL:
			if !ok {
				return ErrSaverTasksWasClosed
			}
			tasksL = nil
			taskResL = task.result
		case result = <-taskResL:
			taskResL = nil
			resultsL = results
		case resultsL <- result:
			resultsL = nil
			tasksL = tasks
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return errors.WithStack(ctx.Err())
}

func saveLoop[TCheckPoint any, TPreSaverData any, TSaverData any](ctx context.Context,
	stor Storage[TCheckPoint, TSaverData],
	prsv PreSaver[TCheckPoint, TPreSaverData, TSaverData],
	results <-chan *parseResult[TPreSaverData],
	maxBatchSize uint, initCheckPoint *TCheckPoint,
) error {
	lastChp := initCheckPoint

	for ctx.Err() == nil {
		b, err := getBatchForSave(ctx, results, maxBatchSize)
		if err != nil {
			return err
		}

		var firstParseErr error
		var psDataList []*TPreSaverData
		for _, r := range b {
			if r.err != nil {
				firstParseErr = r.err
				break
			}

			psDataList = append(psDataList, r.data)
		}

		preSaveChp, sData, err := prsv.PreSaveProcess(ctx, lastChp, psDataList)
		if err != nil {
			return err
		}

		savedChp, err := stor.Save(ctx, preSaveChp, sData)
		if err != nil {
			return err
		}
		lastChp = &savedChp
		if firstParseErr != nil {
			return firstParseErr
		}
	}

	return errors.WithStack(ctx.Err())
}

func getBatchForSave[TPreSaverData any](ctx context.Context,
	results <-chan *parseResult[TPreSaverData],
	maxBatchSize uint,
) ([]*parseResult[TPreSaverData], error) {
	var b []*parseResult[TPreSaverData]

	for len(b) < int(maxBatchSize) {
		var res *parseResult[TPreSaverData]
		if len(b) == 0 {
			select {
			case res = <-results:
			case <-ctx.Done():
				return nil, errors.WithStack(ctx.Err())
			}
		} else {
			select {
			case res = <-results:
			default:
				return b, nil
			}
		}
		b = append(b, res)
	}
	return b, nil
}
