package match

type Sequence interface {
	~string | []Matcher
}
