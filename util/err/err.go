package err

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/go-external-config/go/lang"
)

/*
Recover must be used with defer.

Prints unhandled error stack trace into stderr by default.
*/
func Recover(handler ...func(any)) {
	if err := recover(); err != nil {
		if len(handler) == 0 {
			fmt.Fprintf(os.Stderr, "Unhandled error: %T: %v\n%s", err, err, debug.Stack())
			return
		}
		for _, handler := range handler {
			doHandle(err, handler)
		}
	}
}

func doHandle(err any, handler func(any)) {
	defer Recover()
	handler(err)
}

func StackTrace(skipFrames ...int) []uintptr {
	skip := 2
	if len(skipFrames) > 0 {
		skip = skipFrames[0] + 2
	}
	pcs := make([]uintptr, 32)
	n := runtime.Callers(skip, pcs)
	return append([]uintptr(nil), pcs[:n]...)
}

func PrintStackTrace(err any) string {
	if err == nil {
		return ""
	}
	e, ok := err.(error)
	if !ok {
		return fmt.Sprintf("%T: %v", err, err)
	}
	var b strings.Builder
	for i := 0; e != nil; i++ {
		fmt.Fprintf(&b, lang.If(i == 0, "%T: %v\n", "Caused by: %T: %v\n"), e, e)
		if st, ok := e.(interface{ StackTrace() []uintptr }); ok {
			stack := formatStackTrace(st.StackTrace())
			if stack != "" {
				b.WriteString(stack)
				b.WriteByte('\n')
			}
		}
		e = errors.Unwrap(e)
	}

	return strings.TrimRight(b.String(), "\n")
}

func formatStackTrace(pcs []uintptr) string {
	if len(pcs) == 0 {
		return ""
	}

	var b strings.Builder
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if frame.Function != "runtime.goexit" && frame.Function != "runtime.main" {
			fmt.Fprintf(&b, "    at %s (%s:%d)\n", frame.Function, frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return strings.TrimRight(b.String(), "\n")
}
