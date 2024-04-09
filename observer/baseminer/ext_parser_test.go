package baseminer

import (
	"context"
	"fmt"

	"github.com/anoideaopen/common-component/testshlp"
)

type stubParser struct {
	callHlp     testshlp.CallHlp
	badStubData map[int]error
}

func newStubParser(badStubData map[int]error) *stubParser {
	return &stubParser{
		badStubData: badStubData,
	}
}

func (prsr *stubParser) Parse(_ context.Context, data *stubSrcData) (*stubPreSaverData, error) {
	if err := prsr.callHlp.Call(prsr.Parse); err != nil {
		return nil, err
	}

	if err, ok := prsr.badStubData[data.num]; ok {
		return nil, err
	}

	return &stubPreSaverData{
		num:   data.num,
		data1: fmt.Sprintf("data1_%d", data.num),
		data2: fmt.Sprintf("data2_%d", data.num),
	}, nil
}
