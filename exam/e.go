package exam

import (
	"context"
	"testing"
	"time"

	"github.com/krelinga/go-deep"
)

type Common interface {
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
	Common
	Run(name string, f func(t *testing.T)) bool
}

type E interface {
	Common
	Run(name string, f func(E)) bool
	DeepEnv() deep.Env
}

type eImpl struct {
	TestingT
	deepEnv deep.Env
}

func (t *eImpl) DeepEnv() deep.Env {
	return t.deepEnv
}

func (t *eImpl) Run(name string, f func(E)) bool {
	return t.TestingT.Run(name, func(tt *testing.T) {
		subT := &eImpl{
			TestingT: tt,
			deepEnv:  t.deepEnv,
		}
		f(subT)
	})
}

func New(t TestingT, env deep.Env) E {
	return &eImpl{
		TestingT: t,
		deepEnv:  env,
	}
}
