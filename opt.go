package deep

type Opt interface {
	Update(EnvSetter)
}

type OptFunc func(EnvSetter)

func (f OptFunc) Update(setter EnvSetter) {
	f(setter)
}

type Opts []Opt

func (opts Opts) Update(setter EnvSetter) {
	for _, opt := range opts {
		opt.Update(setter)
	}
}
