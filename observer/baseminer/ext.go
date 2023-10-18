package baseminer

import "context"

type DataParser[TSrcData any, TPreSaverData any] interface {
	Parse(ctx context.Context, data *TSrcData) (*TPreSaverData, error)
}

type PreSaver[TCheckPoint any, TPreSaverData any, TSaverData any] interface {
	PreSaveProcess(ctx context.Context, prevChp *TCheckPoint, data []*TPreSaverData) (TCheckPoint, *TSaverData, error)
}

type Storage[TCheckPoint any, TSaverData any] interface {
	Save(ctx context.Context, chp TCheckPoint, data *TSaverData) (TCheckPoint, error)
	FindCheckPoint(ctx context.Context) (*TCheckPoint, bool, error)
}

type DataSource[TSrcData any] interface {
	SrcData() <-chan *TSrcData
	Close()
}
