package asyncx

import (
	"fmt"
	"runtime/debug"
	"sync"
)

// A RoutineGroup is used to group goroutines together and all wait all goroutines to be done.
type RoutineGroup struct {
	waitGroup sync.WaitGroup
	logFns    []func(string)
}

// NewRoutineGroup returns a RoutineGroup.
func NewRoutineGroup(fns ...func(string)) *RoutineGroup {
	return &RoutineGroup{
		logFns: fns,
	}
}

// Run runs the given fn in RoutineGroup.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *RoutineGroup) Run(fn func()) {
	g.waitGroup.Add(1)

	go func() {
		defer g.waitGroup.Done()
		fn()
	}()
}

// RunSafe runs the given fn in RoutineGroup, and avoid panics.
// Don't reference the variables from outside,
// because outside variables can be changed by other goroutines
func (g *RoutineGroup) RunSafe(fn func()) {
	g.waitGroup.Add(1)

	GoSafe(func() {
		defer g.waitGroup.Done()
		fn()
	}, g.logFns...)
}

// Wait waits all running functions to be done.
func (g *RoutineGroup) Wait() {
	g.waitGroup.Wait()
}

// GoSafe runs the given fn using another goroutine, recovers if fn panics.
func GoSafe(fn func(), logFns ...func(string)) {
	go RunSafe(fn, logFns...)
}

// RunSafe runs the given fn, recovers if fn panics.
func RunSafe(fn func(), logFns ...func(string)) {
	defer func() {
		if p := recover(); p != nil {
			for _, logFn := range logFns {
				logFn(fmt.Sprintf("%+v\n%s", p, debug.Stack()))
			}
		}
	}()

	fn()
}
