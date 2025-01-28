package compute

import "errors"

var (
	ErrStandByParser     = errors.New("stand-by parser")
	ErrEmptyRequest      = errors.New("empty request")
	ErrInvalidArgsNumber = errors.New("invalid args number")
	ErrUnknownCommand    = errors.New("unknown command")
)
