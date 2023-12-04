package state

import "errors"

var (
	ErrBadKey   = errors.New("cannot decode key to address")
	ErrBadValue = errors.New("cannot decode override account")
)
