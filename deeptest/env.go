package deeptest

import (
	"context"
	"testing"
	"time"

	"github.com/krelinga/go-deep"
)

type EnvT interface {
	// Methods from testing.T
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

	// Akin to testing.T's Run method, but with the Env type.
	Run(name string, f func(env Env)) bool
}

type Env interface {
	deep.Env
	EnvT
}

type envImpl struct {
	deep.Env
	*testing.T
}

func (e *envImpl) Run(name string, f func(env Env)) bool {
	return e.T.Run(name, func(t *testing.T) {
		f(WrapEnv(e))
	})
}

func NewEnv(t *testing.T) Env {
	return &envImpl{
		Env: deep.NewEnv(),
		T:   t,
	}
}

type wrappedEnv struct {
	deep.Env
	EnvT
}

func WrapEnv(env Env, opts ...deep.Opt) Env {
	return &wrappedEnv{
		Env:  deep.WrapEnv(env, opts...),
		EnvT: env,
	}
}
