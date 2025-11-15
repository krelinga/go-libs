package match

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-deep/deep"
)

type Slice struct {
	Matchers   []Matcher
	Unordered  bool
	AllowExtra bool
}

func (s Slice) Match(env deep.Env, vals Vals) Result {
	got := deep.AsSlice(vals.Want1())
	if !s.AllowExtra && len(got) != len(s.Matchers) {
		return NewResult(false, fmt.Sprintf("slice length mismatch: got %d, want %d", len(got), len(s.Matchers)))
	}
	if s.Unordered {
		return s.matchUnordered(env, got)
	}
	return s.matchOrdered(env, got)
}

func (s Slice) matchOrdered(env deep.Env, got []reflect.Value) Result {
	matcherIdx, valIdx := 0, 0
	for matcherIdx < len(s.Matchers) && valIdx < len(got) {
		matcher := s.Matchers[matcherIdx]
		val := got[valIdx]
		result := matcher.Match(env, Vals{val})
		if result.Matched() {
			matcherIdx++
		}
		valIdx++
	}
	if matcherIdx < len(s.Matchers) {
		return NewResult(false, fmt.Sprintf("not all matchers matched: matched %d out of %d", matcherIdx, len(s.Matchers)))
	}
	return NewResult(true, "all matchers matched")
}

func (s Slice) matchUnordered(env deep.Env, got []reflect.Value) Result {
	used := make([]bool, len(got))
matchers:
	for i, matcher := range s.Matchers {
		for j, val := range got {
			if used[j] {
				continue
			}
			result := matcher.Match(env, Vals{val})
			if result.Matched() {
				used[j] = true
				continue matchers
			}
		}
		return NewResult(false, fmt.Sprintf("matcher %d did not match any value", i))
	}
	return NewResult(true, "all matchers matched")
}
