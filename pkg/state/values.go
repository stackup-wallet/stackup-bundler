package state

import "errors"

var (
	ErrBadKey   = errors.New("cannot decode override address")
	ErrBadValue = errors.New("cannot decode override account")
)
