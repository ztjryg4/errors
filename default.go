package errors

import "net/http"

// defaultStatusCode 默认状态码，实现了 StatusCode 接口
type defaultStatusCode struct {
	HTTP int
	C    int
	Msg  string
	R    string
}

func (d defaultStatusCode) HTTPCode() int {
	return d.HTTP
}

func (d defaultStatusCode) Code() int {
	return d.C
}

func (d defaultStatusCode) String() string {
	return d.Msg
}

func (d defaultStatusCode) Remark() string {
	return d.R
}

func init() {
	d := defaultStatusCode{
		HTTP: http.StatusInternalServerError,
		C:    0,
		Msg:  "Default status code",
		R:    "Default status code reserved by github.com/ztjryg4/errors",
	}
	u := defaultStatusCode{
		HTTP: http.StatusInternalServerError,
		C:    -1,
		Msg:  "Unknown status code",
		R:    "Unknown status code reserved by github.com/ztjryg4/errors",
	}
	Register(d)
	Register(u)
}
