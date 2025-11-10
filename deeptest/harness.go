package deeptest

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/krelinga/go-deep"
)

type Log struct {
	Error    []string
	Fatal    []string
	Log      []string
	Skip     []string
	Parallel bool
	Fail     bool
	FailNow  bool
	Helper   bool
	SkipNow  bool

	Children map[string]*Log
}

func (l *Log) Find(names ...string) (*Log, bool) {
	if l.Children == nil {
		return nil, false
	}
	cur := l
	for _, name := range names {
		child, ok := cur.Children[name]
		if !ok {
			return nil, false
		}
		cur = child
	}
	return cur, true
}

type Harness struct {
	DeepEnv  deep.Env
	Ctx      context.Context
	Name     string
	Deadline time.Time
}

func (h *Harness) Run(f func(T)) *Log {
	wait := sync.WaitGroup{}
	wait.Add(1)
	log := &Log{}
	go func() {
		defer wait.Done()
		ctx := h.Ctx
		if ctx == nil {
			ctx = context.Background()
		}
		var cancel func()
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
		t := &mockT{
			deepEnv:  h.DeepEnv,
			ctx:      ctx,
			name:     h.Name,
			deadline: h.Deadline,
			log:      log,
		}
		defer t.finish(nil)
		f(t)
	}()

	wait.Wait()
	return log
}

type mockT struct {
	// Initially taken from the harness, but can be modified in sub-tests.
	deepEnv  deep.Env
	ctx      context.Context
	name     string
	deadline time.Time

	log      *Log
	cleanups []func()
	tempDir  string
	done     bool
	mu       sync.Mutex
}

func (t *mockT) finish(result *bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.done = true
	// Run cleanups in reverse order.
	for i := len(t.cleanups) - 1; i >= 0; i-- {
		t.cleanups[i]()
	}
	if result != nil {
		*result = !t.log.Fail
	}
}

func (t *mockT) DeepEnv() deep.Env {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("DeepEnv: T already done")
	}
	return t.deepEnv
}

func (t *mockT) Run(name string, f func(T)) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Run: T already done")
	}

	if _, exists := t.log.Children[name]; exists {
		panic(fmt.Sprintf("Run: sub-test %q already exists", name))
	}
	if t.log.Children == nil {
		t.log.Children = make(map[string]*Log)
	}
	log := &Log{}
	t.log.Children[name] = log
	var result bool
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		defer wait.Done()
		ctx, cancel := context.WithCancel(t.ctx)
		defer cancel()
		subT := &mockT{
			deepEnv:  deep.WrapEnv(t.deepEnv),
			ctx:      ctx,
			name:     fmt.Sprintf("%s/%s", t.name, name),
			deadline: t.deadline,
			log:      log,
		}
		defer subT.finish(&result)
		f(subT)
	}()

	return result
}

func (t *mockT) Chdir(dir string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Chdir: T already done")
	}

	oldDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("Chdir: getwd failed: %w", err))
	}
	err = os.Chdir(dir)
	if err != nil {
		panic(fmt.Errorf("Chdir: chdir to %q failed: %w", dir, err))
	}
	t.cleanupUnlocked(func() {
		err := os.Chdir(oldDir)
		if err != nil {
			panic(fmt.Errorf("Chdir cleanup: chdir to %q failed: %w", oldDir, err))
		}
	})
}

func (t *mockT) cleanupUnlocked(f func()) {
	t.cleanups = append(t.cleanups, f)
}

func (t *mockT) Cleanup(f func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Cleanup: T already done")
	}
	t.cleanupUnlocked(f)
}

func (t *mockT) Context() context.Context {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Context: T already done")
	}

	return t.ctx
}

func (t *mockT) Deadline() (deadline time.Time, ok bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Deadline: T already done")
	}
	return t.deadline, !t.deadline.IsZero()
}

func (t *mockT) Error(args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Error: T already done")
	}

	t.log.Error = append(t.log.Error, fmt.Sprint(args...))
	t.log.Fail = true
}

func (t *mockT) Errorf(format string, args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Errorf: T already done")
	}

	t.log.Error = append(t.log.Error, fmt.Sprintf(format, args...))
	t.log.Fail = true
}

func (t *mockT) Fail() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Fail: T already done")
	}

	t.log.Fail = true
}

func (t *mockT) FailNow() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("FailNow: T already done")
	}

	t.log.FailNow = true
	t.log.Fail = true
	runtime.Goexit()
}

func (t *mockT) Failed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Failed: T already done")
	}

	return t.log.Fail
}

func (t *mockT) Fatal(args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Fatal: T already done")
	}

	t.log.Fatal = append(t.log.Fatal, fmt.Sprint(args...))
	t.log.Fail = true
	t.log.FailNow = true
	runtime.Goexit()
}

func (t *mockT) Fatalf(format string, args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Fatalf: T already done")
	}

	t.log.Fatal = append(t.log.Fatal, fmt.Sprintf(format, args...))
	t.log.Fail = true
	t.log.FailNow = true
	runtime.Goexit()
}

func (t *mockT) Helper() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Helper: T already done")
	}

	t.log.Helper = true
}

func (t *mockT) Log(args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Log: T already done")
	}

	t.log.Log = append(t.log.Log, fmt.Sprint(args...))
}

func (t *mockT) Logf(format string, args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Logf: T already done")
	}

	t.log.Log = append(t.log.Log, fmt.Sprintf(format, args...))
}

func (t *mockT) Name() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Name: T already done")
	}

	return t.name
}

func (t *mockT) Parallel() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Parallel: T already done")
	}

	t.log.Parallel = true
}

func (t *mockT) Setenv(key, value string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Setenv: T already done")
	}

	oldValue, hadOldValue := os.LookupEnv(key)
	err := os.Setenv(key, value)
	if err != nil {
		panic(fmt.Errorf("Setenv: setenv %q=%q failed: %w", key, value, err))
	}
	t.cleanupUnlocked(func() {
		var err error
		if hadOldValue {
			err = os.Setenv(key, oldValue)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			panic(fmt.Errorf("Setenv cleanup: restoring %q failed: %w", key, err))
		}
	})
}

func (t *mockT) Skip(args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Skip: T already done")
	}

	t.log.Skip = append(t.log.Skip, fmt.Sprint(args...))
	t.log.SkipNow = true
	runtime.Goexit()
}

func (t *mockT) Skipf(format string, args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Skipf: T already done")
	}

	t.log.Skip = append(t.log.Skip, fmt.Sprintf(format, args...))
	t.log.SkipNow = true
	runtime.Goexit()
}

func (t *mockT) SkipNow() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("SkipNow: T already done")
	}

	t.log.SkipNow = true
	runtime.Goexit()
}

func (t *mockT) Skipped() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("Skipped: T already done")
	}

	return t.log.SkipNow
}

func (t *mockT) TempDir() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		panic("TempDir: T already done")
	}

	if t.tempDir != "" {
		return t.tempDir
	}

	var err error
	t.tempDir, err = os.MkdirTemp("", "deeptest-*")
	if err != nil {
		panic(fmt.Errorf("TempDir: MkdirTemp failed: %w", err))
	}
	t.cleanupUnlocked(func() {
		err := os.RemoveAll(t.tempDir)
		if err != nil {
			panic(fmt.Errorf("TempDir cleanup: RemoveAll %q failed: %w", t.tempDir, err))
		}
	})

	return t.tempDir
}
