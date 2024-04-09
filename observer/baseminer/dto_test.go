package baseminer

type stubCheckPoint struct {
	num int
	ver int
}

type stubSrcData struct {
	num int
}

type stubPreSaverData struct {
	num   int
	data1 string
	data2 string
}

type stubSaverData struct {
	dataList []*stubPreSaverData
}
