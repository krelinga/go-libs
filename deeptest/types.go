package deeptest

import (
	"context"
	"time"

	"github.com/krelinga/go-deep"
)

// Methods from testing.T (other than Run()) are captured in this interface.
type TestingT interface {
	// TODO: add Attr()
	Chdir(dir string)
	Cleanup(f func())
	Context() context.Context
	Deadline() (deadline time.Time, ok bool)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
	Name() string
	// TODO: add Output()
	Parallel()
	Setenv(key, value string)
	Skip(args ...any)
	SkipNow()
	Skipf(format string, args ...any)
	Skipped() bool
	TempDir() string
}

// Similar to testing.T, but with deep.Env support in the Run() method.
type T interface {
	TestingT
	Run(name string, f func(T, deep.Env)) bool
}