package router

import (
	"errors"
)

type ErrCode interface {
	ErrCode() string
}

// errCode по ошибке определяет код
func errCode(err error) string {
	var errWithCode ErrCode
	if errors.As(err, &errWithCode) {
		return errWithCode.ErrCode()
	}

	switch {
	case errors.Is(err, ErrJsonMarshalResponseData):
		return ErrInternalJsonEncode
	}

	return ErrCodeUnknown
}

const (
	ErrCodeUnknown        = ""
	ErrInternalJsonEncode = "errors.internal.jsonEncode"
	ErrCommonNotFound     = "errors.common.notFound"
)
