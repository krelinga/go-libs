package match

import "errors"

var (
	ErrBadVals  = errors.New("bad vals")
	ErrInternal = errors.New("internal error")
	ErrBadType  = errors.New("bad type")
)
