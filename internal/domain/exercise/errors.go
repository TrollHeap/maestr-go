package exercise

import "errors"

var (
	ErrInvalidStep = errors.New("step index out of range")
	ErrInvalidID   = errors.New("exercise id must be > 0")
)
