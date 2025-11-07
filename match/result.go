package match

type Result interface {
	Matched() bool
	String() string
}

type resultImpl struct {
	matched bool
	string  string
}

func (r *resultImpl) Matched() bool {
	return r.matched
}

func (r *resultImpl) String() string {
	return r.string
}

func NewResult(matched bool, str string) Result {
	return &resultImpl{
		matched: matched,
		string:  str,
	}
}
