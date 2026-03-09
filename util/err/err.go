package err

import (
	"fmt"
	"os"
	"runtime/debug"
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
