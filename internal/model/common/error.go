package common

import "errors"

type Error struct {
	HttpCode  uint
	ErrorCode uint
	err       error
}

func (e Error) IsError() bool {
	return e.err != nil
}

func (e Error) Error() string {
	return e.err.Error()
}

func NewError(code uint, message string) Error {
	return Error{
		HttpCode:  code,
		ErrorCode: code,
		err:       errors.New(message),
	}
}

func NoError() Error {
	return Error{
		HttpCode:  200,
		ErrorCode: 200,
		err:       nil,
	}
}
