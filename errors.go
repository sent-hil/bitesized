package bitesized

import "errors"

var (
	ErrInvalidArg    = errors.New("invalid argument(s)")
	ErrFromAfterTill = errors.New("from date after till")
)
