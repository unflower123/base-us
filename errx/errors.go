package errx

import (
	"fmt"
	"runtime"
)

/**
 * @Author: guyu
 * @Desc:
 * @Date: 2025/4/18
 */
type CodeError struct {
	errCode uint32
	errMsg  string
	file    string
	line    int
}

func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("errcode:%d, errmsg:%s, file:%s, line:%d", e.errCode, e.errMsg, e.file, e.line)
}

func NewErrCodeMsg(errCode uint32, errMsg string) *CodeError {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	funcPC := pcs[0]
	f := runtime.FuncForPC(funcPC)
	var file string
	var line int
	if f != nil {
		file, line = f.FileLine(funcPC)
	}

	return &CodeError{
		errCode: errCode,
		errMsg:  errMsg,
		file:    file,
		line:    line,
	}
}

func NewErrCode(errCode uint32) *CodeError {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	funcPC := pcs[0]
	f := runtime.FuncForPC(funcPC)
	var file string
	var line int
	if f != nil {
		file, line = f.FileLine(funcPC)
	}

	return &CodeError{
		errCode: errCode,
		errMsg:  MapErrMsg(errCode),
		file:    file,
		line:    line,
	}
}
