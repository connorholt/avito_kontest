package models

import (
	"errors"
	"github.com/lib/pq"
)

var (
	ErrNoRecord error = errors.New("models: no matching record found")

	ErrInvalidData = errors.New("models: invalid data")

	ErrDuplicateFeatureTag error = errors.New("models: there's already a banner with a similar feature and a tag")

	ErrDuplicateUserName error = errors.New("storage: duplicate username")

	UniqueViolationErr = pq.ErrorCode("23505")
)

func IsErrorCode(err error, errcode pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}
