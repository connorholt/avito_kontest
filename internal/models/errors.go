package models

import "errors"

var (
	ErrNoRecord error = errors.New("models: no matching record found")

	ErrInvalidData = errors.New("models: invalid data")
)
