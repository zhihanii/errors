package errors

import (
	"fmt"
	"sync"
)

type Coder interface {
	HTTPStatus() int
	String() string
	Reference() string
	Code() int
}

var (
	codes   = map[int]Coder{}
	codeMux = &sync.Mutex{}
)

func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by `github.com/zhihanii/errors` as unknown error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by `github.com/zhihanii/errors` as unknown error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*codeError); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}

	//unknown
	return nil
}

func IsCode(err error, targetCode int) bool {
	if v, ok := err.(*codeError); ok {
		if v.code == targetCode {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, targetCode)
		}

		return false
	}

	return false
}
