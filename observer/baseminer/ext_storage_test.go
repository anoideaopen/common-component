package baseminer

import (
	"context"
	"sync"

	"github.com/atomyze-foundation/common-component/testshlp"
	"github.com/pkg/errors"
)

type stubStorage struct {
	lock    sync.RWMutex
	callHlp testshlp.CallHlp

	data       []*stubPreSaverData
	chp        *stubCheckPoint
	triggerNum int
	triggered  chan struct{}
}

func newStubStorage(triggerNum int) *stubStorage {
	return &stubStorage{
		triggerNum: triggerNum,
		triggered:  make(chan struct{}),
	}
}

func (stor *stubStorage) wasTriggered() <-chan struct{} {
	return stor.triggered
}

func (stor *stubStorage) FindCheckPoint(_ context.Context) (*stubCheckPoint, bool, error) {
	stor.lock.Lock()
	defer stor.lock.Unlock()

	if err := stor.callHlp.Call(stor.FindCheckPoint); err != nil {
		return nil, false, err
	}

	if stor.chp == nil {
		return nil, false, nil
	}

	return stor.chp, true, nil
}

func (stor *stubStorage) Save(_ context.Context, chp stubCheckPoint, data *stubSaverData) (stubCheckPoint, error) {
	stor.lock.Lock()
	defer stor.lock.Unlock()

	if err := stor.callHlp.Call(stor.Save); err != nil {
		return chp, err
	}

	if stor.chp != nil && chp.ver != stor.chp.ver {
		return chp, errors.New("invalid checkpoint version")
	}

	stor.data = append(stor.data, data.dataList...)
	if len(stor.data) > 0 && stor.data[len(stor.data)-1].num == stor.triggerNum {
		close(stor.triggered)
	}

	newChp := stubCheckPoint{
		num: chp.num,
		ver: chp.ver + 1,
	}
	stor.chp = &newChp
	return newChp, nil
}
