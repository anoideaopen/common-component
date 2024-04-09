package baseminer

import (
	"context"

	"github.com/anoideaopen/common-component/testshlp"
)

type stubPreSaver struct {
	callHlp testshlp.CallHlp
	psSuff  string
}

func newPreSaver(psSuff string) *stubPreSaver {
	return &stubPreSaver{
		psSuff: psSuff,
	}
}

func (ps *stubPreSaver) PreSaveProcess(_ context.Context, prevChp *stubCheckPoint, data []*stubPreSaverData) (stubCheckPoint, *stubSaverData, error) {
	if err := ps.callHlp.Call(ps.PreSaveProcess); err != nil {
		return stubCheckPoint{}, nil, err
	}

	chp := stubCheckPoint{}
	if prevChp != nil {
		chp = *prevChp
	}
	if len(data) > 0 {
		chp.num = data[len(data)-1].num
	}

	for _, d := range data {
		d.data1 += ps.psSuff
		d.data2 += ps.psSuff
	}

	return chp, &stubSaverData{dataList: data}, nil
}
