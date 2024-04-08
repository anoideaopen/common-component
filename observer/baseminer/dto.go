package baseminer

type saverTask[TPreSaverData any] struct {
	result <-chan *parseResult[TPreSaverData]
}

type parseResult[TPreSaverData any] struct {
	err  error
	data *TPreSaverData
}

type parserTask[TSrcData any, TPreSaverData any] struct {
	srcData *TSrcData
	result  chan<- *parseResult[TPreSaverData]
}
