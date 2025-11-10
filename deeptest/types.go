package deeptest

import (
	"context"
	"testing"
	"time"

	"github.com/krelinga/go-deep"
)

type CommonT interface {
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

type TestingT interface {
	CommonT
	Run(name string, f func(t *testing.T)) bool
}

// Similar to testing.T, but with deep.Env support.
type T interface {
	CommonT
	Run(name string, f func(T)) bool
	DeepEnv() deep.Env
}

type tImpl struct {
	TestingT
	deepEnv deep.Env
}

func (t *tImpl) DeepEnv() deep.Env {
	return t.deepEnv
}

func (t *tImpl) Run(name string, f func(T)) bool {
	return t.TestingT.Run(name, func(tt *testing.T) {
		subT := &tImpl{
			TestingT: tt,
			deepEnv:  deep.WrapEnv(t.deepEnv),
		}
		f(subT)
	})
}

func NewT(t TestingT, env deep.Env) T {
	return &tImpl{
		TestingT: t,
		deepEnv:  env,
	}
}
