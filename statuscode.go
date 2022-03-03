package errors

import (
	"fmt"
	"sync"
)

// StatusCode 业务状态码
type StatusCode interface {
	// HTTPCode 返回该业务状态对应的 HTTP 状态
	HTTPCode() int

	// Code 返回纯数字的业务状态码
	Code() int

	// String 返回状态码文本信息
	String() string

	// Remark 返回更多备注信息
	Remark() string
}

var (
	statusCodeMap = map[int]StatusCode{}
	mu            = sync.Mutex{}
)

// Register 注册一个业务状态码，重复时覆盖
func Register(code StatusCode) {
	mu.Lock()
	defer mu.Unlock()

	statusCodeMap[code.Code()] = code
}

// MustRegister 注册一个业务状态码，重复时引发 panic
func MustRegister(code StatusCode) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := statusCodeMap[code.Code()]; ok {
		panic(fmt.Sprintf("status code [%d] already exists", code.Code()))
	}

	statusCodeMap[code.Code()] = code
}

// IsStatusCode 检查 err 链中是否包含 StatusCode
func IsStatusCode(err error, code int) bool {
	if v, ok := err.(*withStatusCode); ok {
		if v.code == code {
			return true
		}
		if v.cause != nil {
			return IsStatusCode(v.cause, code)
		}
		return false
	}

	return false
}

// ParseStatusCode 解析任何 error 到 StatusCode
// 不包含业务状态码的 error 将返回保留的 unknownCode
// 空 error 返回 nil
func ParseStatusCode(err error) StatusCode {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withStatusCode); ok {
		if code, ok := statusCodeMap[v.code]; ok {
			return code
		}
	}

	return unknownCode
}
