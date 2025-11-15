package exam_test

import (
	"github.com/krelinga/go-deep/deep"
	"github.com/krelinga/go-deep/exam"
	"github.com/krelinga/go-deep/match"
)

func matchLogOk() match.Matcher {
	return match.Equal(&exam.Log{
		Helper: true,
	})
}

func matchLogError() match.Matcher {
	return match.Pointer(match.Struct{
		Fields: map[deep.Field]match.Matcher{
			deep.NamedField("Error"): match.Len(match.Equal(1)),
		},
		Partial: match.Equal(exam.Log{
			Helper: true,
			Fail:   true,
		}),
	})
}
