package testshlp

import (
	"reflect"
	"runtime"
)

type CallHlp struct {
	callNums  map[string]int
	errorsMap map[string]map[int]error
}

func (ch *CallHlp) AddErrMap(fn interface{}, errMap map[int]error) {
	if ch.callNums == nil {
		ch.callNums = make(map[string]int)
		ch.errorsMap = make(map[string]map[int]error)
	}

	fnName := getFunName(fn)
	ch.callNums[fnName] = 0
	ch.errorsMap[fnName] = errMap
}

func (ch *CallHlp) Call(fn interface{}) error {
	fnName := getFunName(fn)

	num, ok := ch.callNums[fnName]
	if !ok {
		return nil
	}

	ch.callNums[fnName] = num + 1
	em, ok := ch.errorsMap[fnName]
	if !ok {
		return nil
	}
	return em[num]
}

func getFunName(i interface{}) string {
	if i == nil {
		return ""
	}
	if n, ok := i.(string); ok {
		return n
	}

	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
