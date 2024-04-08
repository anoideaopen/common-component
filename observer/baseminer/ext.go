package baseminer

import "context"

// DataParser is a type for data parsers
type DataParser[TSrcData any, TPreSaverData any] interface {
	Parse(ctx context.Context, data *TSrcData) (*TPreSaverData, error)
}

// PreSaver is a type for pre-savers
type PreSaver[TCheckPoint any, TPreSaverData any, TSaverData any] interface {
	PreSaveProcess(ctx context.Context, prevChp *TCheckPoint, data []*TPreSaverData) (TCheckPoint, *TSaverData, error)
}

// Storage is a type for storages
type Storage[TCheckPoint any, TSaverData any] interface {
	Save(ctx context.Context, chp TCheckPoint, data *TSaverData) (TCheckPoint, error)
	FindCheckPoint(ctx context.Context) (*TCheckPoint, bool, error)
}

// DataSource is a type for data sources
type DataSource[TSrcData any] interface {
	SrcData() <-chan *TSrcData
	Close()
}
