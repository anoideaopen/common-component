package baseminer

import (
	"context"

	"github.com/anoideaopen/glog"
	"golang.org/x/sync/errgroup"
)

// MHandler is a type for metrics handler
type MHandler[TCheckPoint any, TSrcData any, TPreSaverData any, TSaverData any] struct {
	log glog.Logger

	parserPoolCountWorkers uint
	saverTasksBufSize      uint
	saverMaxBatchSize      uint

	stor           Storage[TCheckPoint, TSaverData]
	parser         DataParser[TSrcData, TPreSaverData]
	createSrc      func(ctx context.Context, startFromChp *TCheckPoint) (DataSource[TSrcData], error)
	createPreSaver func(ctx context.Context) (PreSaver[TCheckPoint, TPreSaverData, TSaverData], error)
}

// NewMHandler creates a new metrics handler
func NewMHandler[TCheckPoint any, TSrcData any, TPreSaverData any, TSaverData any](ctx context.Context,
	stor Storage[TCheckPoint, TSaverData],
	createSrc func(ctx context.Context, startFromChp *TCheckPoint) (DataSource[TSrcData], error),
	parser DataParser[TSrcData, TPreSaverData],
	createPreSaver func(ctx context.Context) (PreSaver[TCheckPoint, TPreSaverData, TSaverData], error),
	parserPoolCountWorkers, saverTasksBufSize, saverMaxBatchSize uint,
) *MHandler[TCheckPoint, TSrcData, TPreSaverData, TSaverData] {
	log := glog.FromContext(ctx)

	return &MHandler[TCheckPoint, TSrcData, TPreSaverData, TSaverData]{
		log:                    log,
		parserPoolCountWorkers: parserPoolCountWorkers,
		saverTasksBufSize:      saverTasksBufSize,
		saverMaxBatchSize:      saverMaxBatchSize,
		stor:                   stor,
		parser:                 parser,
		createPreSaver:         createPreSaver,
		createSrc:              createSrc,
	}
}

// RunHandler runs metrics handler
func (mh *MHandler[TCheckPoint, TSrcData, TPreSaverData, TSaverData]) RunHandler(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	chp, _, err := mh.stor.FindCheckPoint(ctx)
	if err != nil {
		return err
	}

	prsv, err := mh.createPreSaver(ctx)
	if err != nil {
		return err
	}

	dsrc, err := mh.createSrc(ctx, chp)
	if err != nil {
		return err
	}
	defer dsrc.Close()

	// start processing
	parserTasks := make(chan *parserTask[TSrcData, TPreSaverData], mh.parserPoolCountWorkers)
	eg.Go(func() error {
		return runParserPool(ctx, mh.parser, mh.parserPoolCountWorkers, parserTasks)
	})

	saverTasks := make(chan *saverTask[TPreSaverData], mh.saverTasksBufSize)
	eg.Go(func() error {
		return runSaver(ctx, mh.stor, prsv, saverTasks, mh.saverMaxBatchSize, chp)
	})

	eg.Go(func() error {
		return runProcessor(ctx, dsrc.SrcData(), parserTasks, saverTasks)
	})

	return eg.Wait()
}
