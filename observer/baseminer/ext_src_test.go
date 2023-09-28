package baseminer

import (
	"context"

	"github.com/atomyze-foundation/common-component/testshlp"
)

type stubCrSrc struct {
	callHlp testshlp.CallHlp
	srcData []*stubSrcData
}

func (scr *stubCrSrc) create(ctx context.Context, startFromChp *stubCheckPoint) (DataSource[stubSrcData], error) {
	if err := scr.callHlp.Call(scr.create); err != nil {
		return nil, err
	}

	return newStubSrc(ctx, scr.srcData, startFromChp), nil
}

type stubSrc struct {
	srcData      chan *stubSrcData
	allSrcData   []*stubSrcData
	startFromChp *stubCheckPoint
}

func newStubSrc(ctx context.Context, allSrcData []*stubSrcData, startFromChp *stubCheckPoint) *stubSrc {
	bs := &stubSrc{
		srcData:      make(chan *stubSrcData),
		allSrcData:   allSrcData,
		startFromChp: startFromChp,
	}

	go bs.pushData(ctx)
	return bs
}

func (bs *stubSrc) pushData(ctx context.Context) {
	startFromNum := 0
	if bs.startFromChp != nil {
		startFromNum = bs.startFromChp.num + 1
	}

	for _, data := range bs.allSrcData {
		if ctx.Err() != nil {
			return
		}
		if data.num < startFromNum {
			continue
		}
		select {
		case <-ctx.Done():
			return
		case bs.srcData <- data:
		}
	}
}

func (bs *stubSrc) SrcData() <-chan *stubSrcData {
	return bs.srcData
}

func (bs *stubSrc) Close() {
	close(bs.srcData)
}
