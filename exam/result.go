package exam

type Result interface {
	Log(...any) Result
	Logf(string, ...any) Result
	Ok() bool
	Must()
}

func NewResult(ok bool, e E) Result {
	return &resultImpl{
		ok: ok,
		e:  e,
	}
}

type resultImpl struct {
	ok bool
	e  E
}

func (r *resultImpl) Log(args ...any) Result {
	r.e.Helper()
	if !r.ok {
		r.e.Log(args...)
	}
	return r
}

func (r *resultImpl) Logf(format string, args ...any) Result {
	r.e.Helper()
	if !r.ok {
		r.e.Logf(format, args...)
	}
	return r
}

func (r *resultImpl) Ok() bool {
	return r.ok
}

func (r *resultImpl) Must() {
	r.e.Helper()
	if !r.ok {
		r.e.FailNow()
	}
}
