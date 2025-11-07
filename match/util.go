package match

import "fmt"

func first(primary Matcher, fallbacks ...Matcher) Matcher {
	if primary != nil {
		return primary
	}
	for _, fb := range fallbacks {
		if fb != nil {
			return fb
		}
	}
	panic(fmt.Errorf("%w: no matcher provided", ErrInternal))
}
