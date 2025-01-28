package storage

import "errors"

var (
	ErrStandByEngine     = errors.New("stand-by engine")
	ErrNotFound          = errors.New("entity not found")
	ErrInvalidEntityID   = errors.New("invalid entity id")
	ErrInvalidEntityData = errors.New("invalid entity data")
)
