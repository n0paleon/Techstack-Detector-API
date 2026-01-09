package domain

import "errors"

var (
	ErrInvalidTarget = errors.New("invalid target")
	ErrBlockedTarget = errors.New("target blocked by security policy")
)
