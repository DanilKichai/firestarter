package common

import "errors"

var (
	ErrDataSize           = errors.New("data size is incorrect")
	ErrDataRepresentation = errors.New("data representation is incorrect")
)
