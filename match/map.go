package match

import (
	"fmt"

	"github.com/krelinga/go-deep/deep"
)

type MapEntry struct {
	Key Matcher
	Val Matcher
}

type Map struct {
	Entries    []MapEntry
	AllowExtra bool
}

func (m Map) Match(env deep.Env, vals Vals) Result {
	got := deep.AsMapEntries(vals.Want1())
	used := make([]bool, len(got))
	if !m.AllowExtra && len(got) != len(m.Entries) {
		return NewResult(false, fmt.Sprintf("map length mismatch: got %d, want %d", len(got), len(m.Entries)))
	}
entries:
	for i, entry := range m.Entries {
		for j, kv := range got {
			if used[j] {
				continue
			}
			keyMatcher := first(entry.Key, True())
			keyResult := keyMatcher.Match(env, Vals{kv.Key})
			if !keyResult.Matched() {
				continue
			}
			valMatcher := first(entry.Val, True())
			valResult := valMatcher.Match(env, Vals{kv.Val})
			if !valResult.Matched() {
				continue
			}
			used[j] = true
			continue entries
		}
		return NewResult(false, fmt.Sprintf("entry matcher %d did not match any key-value pair", i))
	}
	return NewResult(true, "all map entries matched")
}
