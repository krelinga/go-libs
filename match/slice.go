package match

import (
	"fmt"

	"github.com/krelinga/go-deep"
)

type Slice struct {
	Matchers  []Matcher
	Unordered bool
}

func (s Slice) Match(env deep.Env, vals Vals) Result {
	got := Want1[[]any](vals)
	if len(got) != len(s.Matchers) {
		return NewResult(false, "slice length mismatch")
	}
	if s.Unordered {
		return s.matchUnordered(env, got)
	}
	return s.matchOrdered(env, got)
}

func (s Slice) matchOrdered(env deep.Env, got []any) Result {
	for i, matcher := range s.Matchers {
		result := matcher.Match(env, NewValsAny(got[i]))
		if !result.Matched() {
			return NewResult(false, fmt.Sprintf("matcher %d did not match: %s", i, result.String()))
		}
	}
	return NewResult(true, "all matchers matched")
}

func (s Slice) matchUnordered(env deep.Env, got []any) Result {
	used := make([]bool, len(got))
matchers:
	for i, matcher := range s.Matchers {
		for j, val := range got {
			if used[j] {
				continue
			}
			result := matcher.Match(env, NewValsAny(val))
			if result.Matched() {
				used[j] = true
				continue matchers
			}
		}
		return NewResult(false, fmt.Sprintf("matcher %d did not match any value", i))
	}
	return NewResult(true, "all matchers matched")
}
