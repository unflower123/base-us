package asyncx

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestRoutineGroupRun test RoutineGroup Run func
func TestRoutineGroupRun(t *testing.T) {
	group := NewRoutineGroup()
	var counter int
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		group.Run(func() {
			mu.Lock()
			counter++
			mu.Unlock()
		})
	}

	group.Wait()

	if counter != 10 {
		t.Errorf("Expected counter to be 10, got %d", counter)
	}
}

// TestRoutineGroupRunSafe test RoutineGroup RunSafe func
func TestRoutineGroupRunSafe(t *testing.T) {
	var logMessage string
	logFn := func(msg string) {
		logMessage = msg
	}

	group := NewRoutineGroup(logFn)
	var counter int
	var mu sync.Mutex

	for i := 0; i < 10; i++ {
		group.RunSafe(func() {
			mu.Lock()
			counter++
			mu.Unlock()
		})
	}

	group.Wait()

	if counter != 10 {
		t.Errorf("Expected counter to be 10, got %d", counter)
	}

	if logMessage != "" {
		t.Errorf("Expected log message to be empty, got %s", logMessage)
	}
}

// TestRoutineGroupRunSafeWithPanic test RunSafe panic Recovery mechanism
func TestRoutineGroupRunSafeWithPanic(t *testing.T) {
	logFn := func(msg string) {
		fmt.Println(msg)
	}

	group := NewRoutineGroup(logFn)

	group.RunSafe(func() {
		panic("test panic")
	})

	time.Sleep(100 * time.Millisecond)

}

// TestGoSafe test GoSafe func
func TestGoSafe(t *testing.T) {
	var logMessage string
	logFn := func(msg string) {
		logMessage = msg
	}

	GoSafe(func() {
		panic("test panic")
	}, logFn)

	time.Sleep(100 * time.Millisecond)

	if logMessage == "" {
		t.Error("Expected log message to be set, got empty string")
	}
}

// TestRunSafe test RunSafe func
func TestRunSafe(t *testing.T) {
	var logMessage string
	logFn := func(msg string) {
		logMessage = msg
	}

	RunSafe(func() {
		panic("test panic")
	}, logFn)

	if logMessage == "" {
		t.Error("Expected log message to be set, got empty string")
	}
}

// TestRoutineGroupWait test whether the Wait method blocks until all tasks are completed
func TestRoutineGroupWait(t *testing.T) {
	group := NewRoutineGroup()
	var counter int
	var mu sync.Mutex

	// start a long-running goroutine
	group.Run(func() {
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		counter++
		mu.Unlock()
	})

	// start a short running goroutine
	group.Run(func() {
		mu.Lock()
		counter++
		mu.Unlock()
	})

	group.Wait()

	if counter != 2 {
		t.Errorf("Expected counter to be 2, got %d", counter)
	}
}

// TestRoutineGroupWithMultipleTasks test concurrent execution of multiple tasks
func TestRoutineGroupWithMultipleTasks(t *testing.T) {
	group := NewRoutineGroup()
	var counter int
	var mu sync.Mutex

	for i := 0; i < 100; i++ {
		group.Run(func() {
			mu.Lock()
			counter++
			mu.Unlock()
		})
	}

	group.Wait()

	if counter != 100 {
		t.Errorf("Expected counter to be 100, got %d", counter)
	}
}
