package deeptest

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/krelinga/go-deep"
)

type MockLog struct {
	Error    []string
	Fatal    []string
	Log      []string
	Skip     []string
	Parallel bool
	Fail     bool
	FailNow  bool
	Helper   bool
	SkipNow  bool
}

type MockEnv interface {
	Env
	MockLog() MockLog
}

type mockEnvImpl struct {
	deep.Env

	// Configuration for the mock testing.T implementation.
	ctx      context.Context
	deadline time.Time
	name     string
	tempDir  string

	// Captured interactions.
	log *MockLog

	// Registered cleanup functions.
	cleanups []func()

	// Mutex to protect concurrent access to the mockEnvImpl.
	mu sync.Mutex
}

func (e *mockEnvImpl) MockLog() MockLog {
	e.mu.Lock()
	defer e.mu.Unlock()

	return MockLog{
		Error:    slices.Clone(e.log.Error),
		Fatal:    slices.Clone(e.log.Fatal),
		Log:      slices.Clone(e.log.Log),
		Skip:     slices.Clone(e.log.Skip),
		Parallel: e.log.Parallel,
		Fail:     e.log.Fail,
		FailNow:  e.log.FailNow,
		Helper:   e.log.Helper,
	}
}

func (e *mockEnvImpl) Chdir(dir string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	oldDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("mockEnvImpl Chdir: unable to get current working directory: %w", err))
	}
	err = os.Chdir(dir)
	if err != nil {
		panic(fmt.Errorf("mockEnvImpl Chdir: unable to change directory to %q: %w", dir, err))
	}

	e.Cleanup(func() {
		err := os.Chdir(oldDir)
		if err != nil {
			panic(fmt.Errorf("mockEnvImpl Chdir cleanup: unable to restore working directory to %q: %w", oldDir, err))
		}
	})
}

func (e *mockEnvImpl) Cleanup(f func()) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.cleanups = append(e.cleanups, f)
}

func (e *mockEnvImpl) Context() context.Context {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ctx != nil {
		return e.ctx
	}
	return context.Background()
}

func (e *mockEnvImpl) Deadline() (deadline time.Time, ok bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.deadline.IsZero() {
		return e.deadline, true
	}
	return time.Time{}, false
}

func (e *mockEnvImpl) Error(args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Error = append(e.log.Error, fmt.Sprint(args...))
	e.log.Fail = true
}

func (e *mockEnvImpl) Errorf(format string, args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Error = append(e.log.Error, fmt.Sprintf(format, args...))
	e.log.Fail = true
}

func (e *mockEnvImpl) Fail() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Fail = true
}

func (e *mockEnvImpl) FailNow() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Fail = true
	e.log.FailNow = true
}

func (e *mockEnvImpl) Failed() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.log.Fail
}

func (e *mockEnvImpl) Fatal(args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Fatal = append(e.log.Fatal, fmt.Sprint(args...))
	e.log.Fail = true
	e.log.FailNow = true
}

func (e *mockEnvImpl) Fatalf(format string, args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Fatal = append(e.log.Fatal, fmt.Sprintf(format, args...))
	e.log.Fail = true
	e.log.FailNow = true
}

func (e *mockEnvImpl) Helper() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Helper = true
}

func (e *mockEnvImpl) Log(args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Log = append(e.log.Log, fmt.Sprint(args...))
}

func (e *mockEnvImpl) Logf(format string, args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Log = append(e.log.Log, fmt.Sprintf(format, args...))
}

func (e *mockEnvImpl) Name() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.name
}

func (e *mockEnvImpl) Parallel() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Parallel = true
}

func (e *mockEnvImpl) Setenv(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	oldEnvValue, hadOldEnvValue := os.LookupEnv(key)
	err := os.Setenv(key, value)
	if err != nil {
		panic(fmt.Errorf("mockEnvImpl Setenv: unable to set environment variable %q: %w", key, err))
	}

	e.Cleanup(func() {
		if hadOldEnvValue {
			err := os.Setenv(key, oldEnvValue)
			if err != nil {
				panic(fmt.Errorf("mockEnvImpl Setenv cleanup: unable to restore environment variable %q: %w", key, err))
			}
		} else {
			err := os.Unsetenv(key)
			if err != nil {
				panic(fmt.Errorf("mockEnvImpl Setenv cleanup: unable to unset environment variable %q: %w", key, err))
			}
		}
	})
}

func (e *mockEnvImpl) Skip(args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Skip = append(e.log.Skip, fmt.Sprint(args...))
	e.log.SkipNow = true
}

func (e *mockEnvImpl) SkipNow() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.SkipNow = true
}

func (e *mockEnvImpl) Skipf(format string, args ...any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.log.Skip = append(e.log.Skip, fmt.Sprintf(format, args...))
	e.log.SkipNow = true
}

func (e *mockEnvImpl) Skipped() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.log.SkipNow
}

func (e *mockEnvImpl) TempDir() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.tempDir
}

func (e *mockEnvImpl) Run(name string, f func(env Env)) bool {
	panic("mockEnvImpl Run: not implemented") // TODO: implement
}

type MockOpt interface {
	updateImpl(e *mockEnvImpl)
}

func MockContext(ctx context.Context) MockOpt {
	return mockContextOpt{ctx: ctx}
}

type mockContextOpt struct {
	ctx context.Context
}

func (o mockContextOpt) updateImpl(e *mockEnvImpl) {
	e.ctx = o.ctx
}

func MockName(name string) MockOpt {
	return mockNameOpt{name: name}
}

type mockNameOpt struct {
	name string
}

func (o mockNameOpt) updateImpl(e *mockEnvImpl) {
	e.name = o.name
}

func MockTempDir(tempDir string) MockOpt {
	return mockTempDirOpt{tempDir: tempDir}
}

type mockTempDirOpt struct {
	tempDir string
}

func (o mockTempDirOpt) updateImpl(e *mockEnvImpl) {
	e.tempDir = o.tempDir
}

func MockDeadline(deadline time.Time) MockOpt {
	return mockDeadlineOpt{deadline: deadline}
}

type mockDeadlineOpt struct {
	deadline time.Time
}

func (o mockDeadlineOpt) updateImpl(e *mockEnvImpl) {
	e.deadline = o.deadline
}

func NewMockEnv(opts ...MockOpt) (mockEnv MockEnv, cleanup func()) {
	e := &mockEnvImpl{
		Env: deep.NewEnv(),
	}
	for _, opt := range opts {
		opt.updateImpl(e)
	}
	return e, func() {
		e.mu.Lock()
		defer e.mu.Unlock()

		// Execute cleanup functions in reverse order.
		for i := len(e.cleanups) - 1; i >= 0; i-- {
			e.cleanups[i]()
		}
	}
}

type Mocker []MockOpt

func (m Mocker) Run(f func(Env)) *MockLog {
	// TODO: implement
	// I think this is a better interface vs. NewMockEnv + defer cleanup()
	// This makes it cleaner for me to start a new goroutine with a mock env
	// and ensure that it finishes before allowing the caller to view the
	// results.
	return nil
}