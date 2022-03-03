package errors

import "net/http"

var (
	unknownCode = ReservedStatusCode{
		HTTP: http.StatusInternalServerError,
		C:    -1,
		Msg:  "An unknown error occurred",
		R:    "Unknown status code reserved by github.com/ztjryg4/errors",
	}
)

// ReservedStatusCode 默认状态码，实现了 StatusCode 接口
type ReservedStatusCode struct {
	HTTP int
	C    int
	Msg  string
	R    string
}

func (r ReservedStatusCode) HTTPCode() int {
	return r.HTTP
}

func (r ReservedStatusCode) Code() int {
	return r.C
}

func (r ReservedStatusCode) String() string {
	return r.Msg
}

func (r ReservedStatusCode) Remark() string {
	return r.R
}

func init() {
	d := ReservedStatusCode{
		HTTP: http.StatusInternalServerError,
		C:    0,
		Msg:  "Default status code",
		R:    "Default status code reserved by github.com/ztjryg4/errors",
	}
	u := unknownCode
	Register(d)
	Register(u)
}
