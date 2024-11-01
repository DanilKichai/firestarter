package unyaml

import "errors"

var (
	ErrTypeRepresentation = errors.New("type representation is incorrect")
	ErrEmptyData          = errors.New("empty data is not allowed")
	ErrEmptyPath          = errors.New("empty path is not allowed")
)
