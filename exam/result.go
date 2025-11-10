package exam

type Result interface {
	Log(...any) Result
	Logf(string, ...any) Result
	Ok() bool
	Must()
}

func NewResult(ok bool, t E) Result {
	return &resultImpl{
		ok: ok,
		e:  t,
	}
}

type resultImpl struct {
	ok bool
	e  E
}

func (r *resultImpl) Log(args ...any) Result {
	if !r.ok {
		r.e.Log(args...)
	}
	return r
}

func (r *resultImpl) Logf(format string, args ...any) Result {
	if !r.ok {
		r.e.Logf(format, args...)
	}
	return r
}

func (r *resultImpl) Ok() bool {
	return r.ok
}

func (r *resultImpl) Must() {
	if !r.ok {
		r.e.FailNow()
	}
}
