// +build go1.7

package errors

import (
	"fmt"
	"testing"

	stderrors "errors"
)

func stdErrors(at, depth int) error {
	if at >= depth {
		return stderrors.New("no error")
	}
	return stdErrors(at+1, depth)
}

func normalErrors(at, depth int) error {
	if at >= depth {
		return New("normal error")
	}
	return normalErrors(at+1, depth)
}

func statusCodeErrors(at, depth int) error {
	if at >= depth {
		return WithStatusCode(0, "status code error")
	}
	return statusCodeErrors(at+1, depth)
}

// GlobalE is an exported global to store the result of benchmark results,
// preventing the compiler from optimising the benchmark functions away.
var GlobalE interface{}

func BenchmarkErrors(b *testing.B) {
	type run struct {
		stack   int
		errType int
	}
	runs := []run{
		{10, 0},
		{10, 1},
		{10, 2},
		{100, 0},
		{100, 1},
		{100, 2},
		{1000, 0},
		{1000, 1},
		{1000, 2},
	}
	for _, r := range runs {
		part := "errors"
		if r.errType == 1 {
			part = "ztjryg4/errors-normal"
		} else if r.errType == 2 {
			part = "ztjryg4/errors-statuscode"
		}
		name := fmt.Sprintf("%s-stack-%d", part, r.stack)
		b.Run(name, func(b *testing.B) {
			var err error
			f := stdErrors
			if r.errType == 1 {
				f = normalErrors
			} else if r.errType == 2 {
				f = statusCodeErrors
			}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = f(0, r.stack)
			}
			b.StopTimer()
			GlobalE = err
		})
	}
}

func BenchmarkStackFormatting(b *testing.B) {
	type run struct {
		stack  int
		format string
	}
	runs := []run{
		{10, "%s"},
		{10, "%v"},
		{10, "%+v"},
		{30, "%s"},
		{30, "%v"},
		{30, "%+v"},
		{60, "%s"},
		{60, "%v"},
		{60, "%+v"},
	}

	var stackStr string
	for _, r := range runs {
		name := fmt.Sprintf("%s-stack-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := normalErrors(0, r.stack)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, err)
			}
			b.StopTimer()
		})
	}

	for _, r := range runs {
		name := fmt.Sprintf("%s-stacktrace-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := normalErrors(0, r.stack)
			st := err.(*fundamental).stack.StackTrace()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, st)
			}
			b.StopTimer()
		})
	}
	GlobalE = stackStr
}

func BenchmarkStatusCodeStackFormatting(b *testing.B) {
	type run struct {
		stack  int
		format string
	}
	runs := []run{
		{10, "%s"},
		{10, "%v"},
		{10, "%+v"},
		{10, "%#v"},
		{10, "%#+v"},
		{30, "%s"},
		{30, "%v"},
		{30, "%+v"},
		{30, "%#v"},
		{30, "%#+v"},
		{60, "%s"},
		{60, "%v"},
		{60, "%+v"},
		{60, "%#v"},
		{60, "%#+v"},
	}

	var stackStr string
	for _, r := range runs {
		name := fmt.Sprintf("%s-stack-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := statusCodeErrors(0, r.stack)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, err)
			}
			b.StopTimer()
		})
	}

	for _, r := range runs {
		name := fmt.Sprintf("%s-stacktrace-%d", r.format, r.stack)
		b.Run(name, func(b *testing.B) {
			err := statusCodeErrors(0, r.stack)
			st := err.(*withStatusCode).stack.StackTrace()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				stackStr = fmt.Sprintf(r.format, st)
			}
			b.StopTimer()
		})
	}
	GlobalE = stackStr
}
