package concurrent

import (
	"sync"
	"time"
)

func Synchronized(mutex *sync.Mutex, operation func()) {
	mutex.Lock()
	defer mutex.Unlock()
	operation()
}

func SynchronizedWithCondition(mutex *sync.Mutex, cond *sync.Cond, operation func(cond *sync.Cond)) {
	mutex.Lock()
	defer mutex.Unlock()
	operation(cond)
}

// CondWaitTimeout waits on condition until it is signaled or timeout elapses.
// It returns true if the timeout timer did NOT fire before we stopped it
// (i.e. wake was likely due to some other Signal/Broadcast),
// and false if the timeout timer has been fired.
func WaitWithTimeout(cond *sync.Cond, timeout time.Duration) bool {
	timer := time.AfterFunc(timeout, func() {
		cond.Broadcast()
	})
	cond.Wait()
	return timer.Stop()
}
