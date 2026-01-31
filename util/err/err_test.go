package err_test

import (
	"testing"

	"github.com/go-external-config/go/util/err"
)

func TestRecover_SwallowsPanic(t *testing.T) {
	didReachEnd := false

	func() {
		defer err.Recover()
		panic("boom")
		didReachEnd = true // unreachable
	}()

	if didReachEnd {
		t.Fatalf("execution should not continue after panic point")
	}
}

func TestRecover_WithHandler_CallsHandler(t *testing.T) {
	var called bool
	var got any

	func() {
		defer err.Recover(func(v any) {
			called = true
			got = v
		})
		panic(123)
	}()

	if !called {
		t.Fatalf("handler was not called")
	}
	if got != 123 {
		t.Fatalf("expected handler to receive 123, got %#v", got)
	}
}

func TestRecover_HandlerPanic_IsRecovered(t *testing.T) {
	var secondCalled bool

	func() {
		defer err.Recover(
			func(any) {
				panic("handler failed")
			},
			func(any) {
				secondCalled = true
			},
		)
		panic("original")
	}()

	// if handler panic was not recovered, the test would crash
	if !secondCalled {
		t.Fatalf("expected second handler to be called even if first panics")
	}
}
